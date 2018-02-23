package curator

import (
	"sync"

	"github.com/samuel/go-zookeeper/zk"
)

type Watcher struct {
	f func(event zk.Event)
}

func NewWatcher(f func(zk.Event)) *Watcher {
	return &Watcher{f}
}

func (w *Watcher) Fire(event zk.Event) {
	w.f(event)
}

type watcherManager struct {
	mutex    sync.Mutex
	watchers map[*Watcher]*Watcher
}

func newWatcherManager() *watcherManager {
	return &watcherManager{
		watchers: make(map[*Watcher]*Watcher),
	}
}

func (w *watcherManager) AddWatcher(watcher *Watcher) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.watchers[watcher] = watcher
}

func (w *watcherManager) DelWatcher(watcher *Watcher) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	delete(w.watchers, watcher)
}

func (w *watcherManager) Fire(event zk.Event) {
	m := make(map[*Watcher]*Watcher)

	w.mutex.Lock()
	for k, v := range w.watchers {
		m[k] = v
	}
	w.mutex.Unlock()

	for _, v := range m {
		v.Fire(event)
	}
}
