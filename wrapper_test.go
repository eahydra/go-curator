package curator

import (
	"reflect"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
)

const aclTestNode = "/test/acltest"

func TestZookeeperClient_GetACL(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to newZookeeperClient, err:", err)
	}
	defer client.Close()

	expectedACL := zk.WorldACL(zk.PermAll)
	CreateAll(client, aclTestNode, nil, zk.FlagEphemeral, expectedACL)
	defer DeleteAll(client, aclTestNode)
	acl, _, err := client.GetACL(aclTestNode)
	if err != nil {
		t.Fatal("failed to client.GetACL, err:", err)
	}
	if !reflect.DeepEqual(acl, expectedACL) {
		t.Fatal("unexpected acl")
	}
}

func TestZookeeperClient_SetACL(t *testing.T) {
	client, err := newZooKeeperClient()
	if err != nil {
		t.Fatal("failed to newZookeeperClient, err:", err)
	}
	defer client.Close()

	CreateAll(client, aclTestNode, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	defer DeleteAll(client, aclTestNode)
	expectedACL := zk.WorldACL(zk.PermRead)
	_, err = client.SetACL(aclTestNode, expectedACL, -1)
	if err != nil {
		t.Fatal("failed to client.SetACL, err:", err)
	}

	acl, _, err := client.GetACL(aclTestNode)
	if err != nil {
		t.Fatal("failed to client.GetACL, err:", err)
	}

	if !reflect.DeepEqual(acl, expectedACL) {
		t.Fatal("unexpected acl")
	}
}
