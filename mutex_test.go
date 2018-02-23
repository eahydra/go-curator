package curator

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func TestMutexSortChildren(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	count := 100
	children := make([]string, 0, count)
	prefix := "lock-"
	for i := 0; i <= count; i++ {
		children = append(children, fmt.Sprintf("%s%d", prefix, rand.Int()))
	}
	sort.Sort(&mutexSortChildren{prefix, children})
	if !sort.IsSorted(&mutexSortChildren{prefix, children}) {
		t.Fatal("sorted failed")
	}
}

func TestMutex(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	basePath := "/test/mutex_test"

	mutex := NewMutex(client, basePath, zk.WorldACL(zk.PermAll))
	if err := mutex.Acquire(); err != nil {
		t.Fatal(err)
	}
	t.Log(mutex.lockPath)
	mutex.Release()

	ch := make(chan error, 1)
	participantMutex := NewMutex(client, basePath, zk.WorldACL(zk.PermAll))
	go func() {
		if err := participantMutex.Acquire(); err != nil {
			ch <- err
			return
		}

		ch <- nil
		time.Sleep(1 * time.Second)
		participantMutex.Release()
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("participant cannot acquire mutex")
	case err := <-ch:
		if err != nil {
			t.Fatal("participant failed to acquire mutex, err:", err)
		}
	}

	if err := mutex.Acquire(); err != nil {
		t.Fatal(err)
	}

	mutex.Release()
}

func TestMutex_Release(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	basePath := "/test/mutex_test"

	mutex := NewMutex(client, basePath, zk.WorldACL(zk.PermAll))
	err = mutex.Release()
	if err != zk.ErrNotLocked {
		t.Fatal("unexpected err on mutex.Release, err:", err)
	}

	err = mutex.Acquire()
	if err != nil {
		t.Fatal("failed to mutex.Acquire, err:", err)
	}
	err = mutex.Release()
	if err != nil {
		t.Fatal("failed to mutex.Release, err:", err)
	}
}
