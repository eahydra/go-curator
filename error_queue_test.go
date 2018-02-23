package curator

import (
	"errors"
	"testing"
)

func TestErrorQueue_Push(t *testing.T) {
	q := newErrorQueue(1)
	q.Push(nil)
	q.Push(nil)
	if q.Len() != 1 {
		t.Fatal("failed to errorQueue.Push")
	}
}

func TestErrorQueue_Pop(t *testing.T) {
	q := newErrorQueue(1)
	q.Push(nil)
	err := errors.New("error")
	q.Push(err)
	e := q.Pop()
	if e != err {
		t.Fatal("failed to q.Pop")
	}
}
