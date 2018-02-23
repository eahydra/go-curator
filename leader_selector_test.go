package curator

import (
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type MockLeaderSelectorListener struct {
	done chan struct{}
}

func (m *MockLeaderSelectorListener) TakeLeaderShip(client *ZookeeperClient, cancel <-chan struct{}) error {
	close(m.done)
	select {
	case <-time.After(1 * time.Second):
	case <-cancel:
	}
	return nil
}

func TestLeaderSelector(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	listener := &MockLeaderSelectorListener{
		done: make(chan struct{}, 1),
	}
	leaderSelector := NewLeaderSelector(client, "/test/leader-test", listener, zk.WorldACL(zk.PermAll))
	leaderSelector.Start()
	defer leaderSelector.Close()

	select {
	case <-listener.done:
		if !leaderSelector.HasLeaderShip() {
			t.Fatal("unexpected leadership")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("leader failed")
	}
}
