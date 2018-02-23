package curator

import (
	"errors"
	"path"
	"sync"
	"sync/atomic"

	"github.com/samuel/go-zookeeper/zk"
)

type ChildrenCacheEvent struct {
	ChildNode string
	Data      []byte
	Stat      *zk.Stat
	Type      ChildrenCacheEventType
}

type ChildrenCacheEventType int

const (
	ChildrenCacheAdd    ChildrenCacheEventType = 1
	ChildrenCacheUpdate ChildrenCacheEventType = 2
	ChildrenCacheDel    ChildrenCacheEventType = 3
)

func (c ChildrenCacheEventType) String() string {
	switch c {
	case ChildrenCacheAdd:
		return "ChildrenCacheAdd"
	case ChildrenCacheDel:
		return "ChildrenCacheDel"
	case ChildrenCacheUpdate:
		return "ChildrenCacheUpdate"
	}
	return "unknown"
}

type OnChildrenCacheChange func(event ChildrenCacheEvent)

type childContext struct {
	done chan struct{}
	wg   sync.WaitGroup
	lock sync.RWMutex
	stat *zk.Stat
	data []byte
}

func newChildContext() *childContext {
	return &childContext{
		done: make(chan struct{}, 1),
	}
}

func (c *childContext) getDataAndStat() ([]byte, *zk.Stat) {
	c.lock.RLock()
	data, stat := c.data, c.stat
	c.lock.RUnlock()
	return data, stat
}

func (c *childContext) update(data []byte, stat *zk.Stat) {
	c.lock.Lock()
	c.data, c.stat = data, stat
	c.lock.Unlock()
}

type ChildrenCache struct {
	start    int32
	client   *ZookeeperClient
	node     string
	callback OnChildrenCacheChange
	mutex    *sync.RWMutex
	cache    map[string]*childContext
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewChildrenCache(client *ZookeeperClient, node string, callback OnChildrenCacheChange) *ChildrenCache {
	return &ChildrenCache{
		client:   client,
		node:     node,
		callback: callback,
		mutex:    new(sync.RWMutex),
		cache:    make(map[string]*childContext),
	}
}

func (w *ChildrenCache) Start() error {
	if !atomic.CompareAndSwapInt32(&w.start, 0, 1) {
		return errors.New("curator: ChildrenCache already started")
	}

	exist, _, err := w.client.Exists(w.node)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("curator: node is not exist, node: " + w.node)
	}

	w.quit = make(chan struct{}, 1)
	w.wg.Add(1)
	go w.watchChildren()
	return nil
}

func (w *ChildrenCache) Close() error {
	if !atomic.CompareAndSwapInt32(&w.start, 1, 0) {
		return errors.New("curator: ChildrenCache already closed")
	}

	close(w.quit)
	w.wg.Wait()
	return nil
}

func (w *ChildrenCache) Get(child string) (data []byte, stat *zk.Stat, ok bool) {
	w.mutex.RLock()
	ctx := w.cache[child]
	w.mutex.RUnlock()
	if ctx != nil {
		ok = true
		data, stat = ctx.getDataAndStat()
	}
	return
}

func (w *ChildrenCache) watchChildren() {
	Log.Infoln("curator: start ChildrenCache.watchChildren", w.node)
	defer func() {
		Log.Infoln("curator.zk: stop ChildrenCache.watchChildren", w.node)
		w.mutex.Lock()
		for _, ctx := range w.cache {
			close(ctx.done)
			ctx.wg.Wait()
		}
		w.mutex.Unlock()

		w.wg.Done()
	}()

	for {
		children, _, eventChan, err := w.client.ChildrenW(w.node)
		if err != nil {
			Log.Errorln("curator: failed to client.ChildrenW, node:", w.node, "err:", err)
			return
		}

		w.mutex.Lock()
		oldCache := make(map[string]*childContext)
		newCache := make(map[string]*childContext)
		for _, child := range children {
			if v, ok := w.cache[child]; ok {
				oldCache[child] = v
				delete(w.cache, child)
				continue
			}
			newCache[child] = newChildContext()
		}

		for _, ctx := range w.cache {
			close(ctx.done)
			ctx.wg.Wait()
		}
		w.cache = oldCache

		for child, ctx := range newCache {
			w.cache[child] = ctx
			ctx.wg.Add(1)
			go w.watchChild(child, ctx)
		}
		w.mutex.Unlock()

		select {
		case <-w.quit:
			return
		case event, ok := <-eventChan:
			if !ok {
				return
			}

			if event.Type == zk.EventNodeDeleted {
				Log.Warnln("curator: ChildrenCache find node had been deleted, node:", w.node)
				return
			}
		}
	}
}

func (w *ChildrenCache) watchChild(child string, ctx *childContext) {
	childPath := path.Join(w.node, child)
	Log.Infoln("curator: start ChildrenCache.watchChild", childPath)

	defer func() {
		Log.Infoln("curator: stop ChildrenCache.watchChild", childPath)
		ctx.wg.Done()
	}()

	eventType := ChildrenCacheAdd
	for {
		data, stat, eventChan, err := w.client.GetW(childPath)
		if err != nil {
			Log.Warnln("curator: failed to GetW, childPath:", childPath, "err:", err)
			return
		}

		if eventType == ChildrenCacheAdd || eventType == ChildrenCacheUpdate {
			ctx.update(data, stat)
		}

		w.notify(ChildrenCacheEvent{ChildNode: childPath, Data: data, Stat: stat, Type: eventType})

		if eventType == ChildrenCacheAdd {
			eventType = ChildrenCacheUpdate
		}

		select {
		case <-w.quit:
			return
		case <-ctx.done:
			return
		case event, ok := <-eventChan:
			if !ok {
				return
			}

			if event.Type == zk.EventNodeDeleted {
				w.notify(ChildrenCacheEvent{ChildNode: childPath, Type: ChildrenCacheDel})
				return
			}
		}
	}
}

func (w *ChildrenCache) notify(event ChildrenCacheEvent) {
	Log.Infoln("curator: ChildrenCache node:", event.ChildNode, "event:", event.Type)

	if w.callback != nil {
		w.callback(event)
	}
}
