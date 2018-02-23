package curator

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/samuel/go-zookeeper/zk"
)

type LeaderSelectorListener interface {
	TakeLeaderShip(client *ZookeeperClient, cancel <-chan struct{}) error
}

type LeaderSelector struct {
	client     *ZookeeperClient
	mutex      *Mutex
	start      int32
	leaderShip int32
	listener   LeaderSelectorListener
	done       chan struct{}
	wg         *sync.WaitGroup
}

func NewLeaderSelector(client *ZookeeperClient, basePath string, listener LeaderSelectorListener, aclv []zk.ACL) *LeaderSelector {
	return &LeaderSelector{
		client:   client,
		mutex:    NewMutex(client, basePath, aclv),
		listener: listener,
		wg:       new(sync.WaitGroup),
	}
}

func (l *LeaderSelector) Start() error {
	if !atomic.CompareAndSwapInt32(&l.start, 0, 1) {
		return errors.New("curator: LeaderSelector already started")
	}
	l.done = make(chan struct{}, 1)
	l.wg.Add(1)
	go l.workLoop()
	return nil
}

func (l *LeaderSelector) Close() error {
	if atomic.CompareAndSwapInt32(&l.start, 1, 0) {
		close(l.done)
		l.wg.Wait()
	}
	return nil
}

func (l *LeaderSelector) HasLeaderShip() bool {
	return atomic.LoadInt32(&l.leaderShip) == 1
}

func (l *LeaderSelector) work() {
	if err := l.mutex.Acquire(); err != nil {
		return
	}

	atomic.StoreInt32(&l.leaderShip, 1)

	connectionCh := make(chan zk.Event, 1)
	watcher := NewWatcher(func(ev zk.Event) {
		if ev.State == zk.StateExpired || ev.State == zk.StateDisconnected {
			select {
			case connectionCh <- ev:
			default:
			}
		}
	})
	l.client.AddWatcher(watcher)

	cancel := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		atomic.StoreInt32(&l.leaderShip, 0)
		l.mutex.Release()
		l.client.DelWatcher(watcher)
		close(cancel)
		wg.Wait()
	}()

	runErrCh := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErrCh <- l.listener.TakeLeaderShip(l.client, cancel)
	}()

	select {
	case <-runErrCh:
	case <-connectionCh:
	case <-l.done:
	}
}

func (l *LeaderSelector) workLoop() {
	defer l.wg.Done()
	for atomic.LoadInt32(&l.start) == 1 {
		l.work()
	}
}
