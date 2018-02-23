package curator

import (
	"strconv"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func TestCreeteAll(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to NewZooKeeperClient, err:", err)
	}
	defer client.Close()

	nodePath := "/zookeeper/" + strconv.FormatInt(time.Now().UnixNano(), 10)
	value := "Hello CreateAll"
	pathCreated, err := CreateAll(client, nodePath, []byte(value), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		t.Fatal("failed to CreateAll, err:", err)
	}
	t.Log("CreateAll", pathCreated, "succeeded")

	client.Get(nodePath)
	data, _, err := client.Get(nodePath)
	if err != nil {
		t.Fatal("failed to client.Get, err:", err)
	}
	if value != string(data) {
		t.Fatal("value is different")
	}

	if _, err := CreateAll(client, nodePath, []byte(value), zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); err != zk.ErrNodeExists {
		t.Fatal("unexpected error:", err)
	}

	DeleteAll(client, nodePath)
}
