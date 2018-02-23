package curator

import (
	"container/list"
	"sync"
)

type errorQueue struct {
	mutex sync.Mutex
	max   int
	l     *list.List
}

func newErrorQueue(max int) *errorQueue {
	return &errorQueue{
		max: max,
		l:   list.New(),
	}
}

func (e *errorQueue) Push(err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for e.l.Len() >= e.max {
		if v := e.l.Front(); v != nil {
			e.l.Remove(v)
		}
	}

	e.l.PushBack(err)
}

func (e *errorQueue) Pop() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var err error
	v := e.l.Front()
	if v != nil {
		err = v.Value.(error)
		e.l.Remove(v)
	}
	return err
}

func (e *errorQueue) Len() int {
	e.mutex.Lock()
	n := e.l.Len()
	e.mutex.Unlock()
	return n
}
