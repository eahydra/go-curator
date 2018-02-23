package curator

import (
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var testServers = "192.168.191.28:2181,192.168.191.29:2181,192.168.191.30:2181"

func newZooKeeperClient() (*ZookeeperClient, error) {
	client, err := NewZookeeperClientBuidler().
		WithZookeeperFactory(DefaultZookeeperFactory).
		WithEnsembleProvider(NewFixedEnsembleProvider(testServers)).
		WithRetryPolicy(NewRetryForever(500 * time.Millisecond)).
		WithSessionTimeout(3 * time.Second).
		WithConnectionTimeout(1 * time.Second).
		WithCanBeReadOnly(true).
		Build()
	if err != nil {
		return nil, err
	}

	if err := client.Start(); err != nil {
		return nil, err
	}

	return client, nil
}

func TestZKClient(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to NewZooKeeperClient, err:", err)
	}
	defer client.Close()

	const testNodePath = "/zookeeper/test"

	pathCreated, err := client.Create(testNodePath, []byte("Hello ZooKeeper"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		if err != zk.ErrNodeExists {
			t.Fatal("failed to client.Create, err:", err)
		}
		t.Log("client.Create", testNodePath, "exists")
		pathCreated = testNodePath
	}
	t.Log("client.Create return", pathCreated)

	var waiter sync.WaitGroup
	waiter.Add(1)

	data, _, eventChan, err := client.GetW(pathCreated)
	if err != nil {
		t.Log("failed to client.GetW, err:", err)
		t.FailNow()
	}

	t.Log("get from", pathCreated, "data:", string(data))

	go func() {
		defer waiter.Done()

		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					t.Log("eventChan is closed")
					return
				}
				t.Log("event:", event)
				if event.Type == zk.EventNodeDataChanged {
					t.Log("data is changed")
					data, _, eventChan, err = client.GetW(pathCreated)
					if err != nil {
						t.Log("failed to client.GetW, err:", err)
					} else {
						t.Log("new data:", string(data))
					}
					return
				}
			}
		}
	}()

	waiter.Add(1)
	go func() {
		defer waiter.Done()

		_, err := client.Set(pathCreated, []byte("Hello GetW, I changes the values."), -1)
		if err != nil {
			t.Log("failed to client.Set, err:", err)
			return
		}

		t.Log("client.Set succeeded")
	}()

	waiter.Wait()
}

func TestClientClose(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to NewZooKeeperClient, err:", err)
	}
	client.Close()
	t.Log("client closed.")

	t.Log("try to client.GetW")
	_, _, _, err = client.GetW("/zookeeper")
	if err != nil {
		t.Log("client.GetW return err, it's succeeded. err:", err)
		return
	}
}

func TestClientConnectChange(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to NewZooKeeperClient, err:", err)
	}
	defer client.Close()

	done := make(chan zk.State, 1)
	watcher := NewWatcher(func(event zk.Event) {
		select {
		case done <- event.State:
		default:
		}
	})
	client.AddWatcher(watcher)
	client.Close()

	select {
	case <-done:
		t.Log("connecting is changed")
	case <-time.After(1 * time.Second):
		t.Fatal("cannot get connection state")
	}
}

func TestClientChildren(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to NewZooKeeperClient, err:", err)
	}
	defer client.Close()

	value := []byte("TestClientChildren")

	const testNodePath = "/zookeeper/abc"
	_, err = client.Create(testNodePath, value, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		if err != zk.ErrNodeExists {
			t.Fatal("failed to client.Create", testNodePath, "err:", err)
		}
	}

	for i := 0; i < 10; i++ {
		p := path.Join(testNodePath, strconv.Itoa(i))
		_, err := client.Create(p, value, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			if err != zk.ErrNodeExists {
				t.Fatal("failed to client.Create child, path:", p, "err:", err)
			}
		}
	}

	children, _, err := client.Children(testNodePath)
	if err != nil {
		t.Fatal("failed to client.Children, err:", err)
	}
	for _, c := range children {
		t.Log("child:", c)
	}
}

func TestClientChildrenW(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to NewZooKeeperClient, err:", err)
	}
	defer client.Close()

	value := []byte("TestClientChildren")
	const root = "/zookeeper/TestClientChildrenW"
	if _, err := client.Create(root, value, 0, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			t.Fatal("failed to client.Create, err:", err)
		}
	}

	if _, err := client.Create(path.Join(root, "/test"), value, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			t.Fatal("failed to client.Create, err:", err)
		}
	}

	children, _, watch, err := client.ChildrenW(root)
	if err != nil {
		t.Fatal("failed to client.ChildrenW, err:", err)
	}

	for _, child := range children {
		t.Log("child:", child)
	}

	p := path.Join(root, "new")
	_, err = client.Create(p, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		t.Fatal("failed to client.Create with ", p, "err:", err)
	}

	select {
	case event, ok := <-watch:
		if !ok {
			return
		}
		t.Log(event)
		return
	case <-time.After(10 * time.Second):
		return
	}
}

func TestNoServer(t *testing.T) {
	ensemble := NewFixedEnsembleProvider("127.0.0.1:2181")
	sessionTimeout := 3 * time.Second
	connectionTimeout := 1 * time.Second
	retryPolicy := NewRetryNTimes(3, 100*time.Millisecond)

	client, err := NewZookeeperClient(DefaultZookeeperFactory, ensemble, sessionTimeout, connectionTimeout, retryPolicy, true)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Start(); err != nil {
		t.Fatal(err)
	}

	count := 0
	err = CallWithRetryLoop(client, func() error {
		count++
		_, _, _, err := client.GetConn().GetW("/zookeeer")
		return err
	})
	if count != 4 {
		t.Fatal("unexpected retry count", count, err)
	}
}
