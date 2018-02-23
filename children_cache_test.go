package curator

import (
	"bytes"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const childrenCacheNode = "/test/childrenCache"

func TestChildrenCache_Start(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to newZookeeperClient, err:", err)
	}
	defer client.Close()

	var expectedNode []string
	for i := 0; i < 3; i++ {
		node := path.Join(childrenCacheNode, strconv.Itoa(i))
		expectedNode = append(expectedNode, node)
		_, err = CreateAll(client, node, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			t.Fatal("failed to Create, err:", err)
		}
	}
	defer DeleteAll(client, childrenCacheNode)

	var nodes []string
	c := make(chan string, len(expectedNode))
	cache := NewChildrenCache(client, childrenCacheNode, func(event ChildrenCacheEvent) {
		c <- event.ChildNode
	})
	err = cache.Start()
	if err != nil {
		t.Fatal("failed to cache.Start, err:", err)
	}
	defer cache.Close()

	err = cache.Start()
	if err == nil {
		t.Fatal("unexpected Start result")
	}

	deadline := time.After(1 * time.Second)
recvLoop:
	for i := 0; i < len(expectedNode); i++ {
		select {
		case childNode := <-c:
			nodes = append(nodes, childNode)
		case <-deadline:
			break recvLoop
		}
	}

	if len(expectedNode) != len(nodes) {
		t.Fatal("unexpected nodes, expected:", expectedNode, "but:", nodes)
	}
	for _, v := range expectedNode {
		found := false
		for _, vv := range nodes {
			if v == vv {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("unexpected result, not found target:", v)
		}
	}
}

func TestChildrenCache_Close(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to newZookeeperClient, err:", err)
	}
	defer client.Close()

	CreateAll(client, childrenCacheNode, nil, 0, zk.WorldACL(zk.PermAll))

	cache := NewChildrenCache(client, childrenCacheNode, nil)
	err = cache.Start()
	if err != nil {
		t.Fatal("failed to cache.Start, err:", err)
	}
	err = cache.Close()
	if err != nil {
		t.Fatal("failed to cache.Close, err:", err)
	}

	err = cache.Close()
	if err == nil {
		t.Fatal("failed to cache.Close twice, must return error")
	}

	err = cache.Start()
	if err != nil {
		t.Fatal("failed to cache.Start, err:", err)
	}
	cache.Close()
}

func TestChildrenCache_Get(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to newZookeeperClient, err:", err)
	}
	defer client.Close()

	CreateAll(client, childrenCacheNode, nil, 0, zk.WorldACL(zk.PermAll))

	cache := NewChildrenCache(client, childrenCacheNode, nil)
	err = cache.Start()
	if err != nil {
		t.Fatal("failed to cache.Start, err:", err)
	}
	defer cache.Close()

	expectData := []byte("hello childrencache")
	CreateAll(client, path.Join(childrenCacheNode, "123"), expectData, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	time.Sleep(1 * time.Second)
	defer DeleteAll(client, childrenCacheNode)

	data, stat, ok := cache.Get("123")
	if !ok {
		t.Fatal("failed to cache.Get")
	}
	if stat == nil {
		t.Fatal("stat must not be nil")
	}
	if len(data) != len(expectData) {
		t.Fatal("unexpected value")
	}
	if bytes.Compare(data, expectData) != 0 {
		t.Fatal("unexpected value")
	}
}

func TestChildCache_StartWithInvalidNode(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to newZookeeperClient, err:", err)
	}
	defer client.Close()

	CreateAll(client, childrenCacheNode, nil, 0, zk.WorldACL(zk.PermAll))
	expectedData := []byte("hello children cache")
	c := make(chan []byte, 1)
	cache := NewChildrenCache(client, childrenCacheNode, func(event ChildrenCacheEvent) {
		c <- event.Data
	})
	err = cache.Start()
	if err != nil {
		t.Fatal("failed to cache.Start")
	}
	defer cache.Close()

	CreateAll(client, path.Join(childrenCacheNode, "update"), expectedData, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	deadline := time.After(1 * time.Second)
	select {
	case <-deadline:
		t.Fatal("deadline")
	case d := <-c:
		if bytes.Compare(expectedData, d) != 0 {
			t.Fatal("unexpected value")
		}
	}
	d, _, ok := cache.Get("update")
	if !ok {
		t.Fatal("unexpected Get result")
	}
	if bytes.Compare(expectedData, d) != 0 {
		t.Fatal("unexpected value")
	}
}

func TestChildrenCache_Update(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

}
