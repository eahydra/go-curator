package curator

import (
	"github.com/samuel/go-zookeeper/zk"
)

type Conn interface {
	Close()

	Get(path string) (data []byte, stat *zk.Stat, err error)
	GetW(path string) (data []byte, stat *zk.Stat, watch <-chan zk.Event, err error)

	Children(path string) (children []string, stat *zk.Stat, err error)
	ChildrenW(path string) (children []string, stat *zk.Stat, watch <-chan zk.Event, err error)

	Exists(path string) (exist bool, stat *zk.Stat, err error)
	ExistsW(path string) (exist bool, stat *zk.Stat, watch <-chan zk.Event, err error)

	Create(path string, value []byte, flags int32, acl []zk.ACL) (pathCreated string, err error)
	CreateProtectedEphemeralSequential(path string, value []byte, aclv []zk.ACL) (pathCreated string, err error)

	Set(path string, value []byte, version int32) (stat *zk.Stat, err error)

	Delete(path string, version int32) (err error)

	GetACL(path string) ([]zk.ACL, *zk.Stat, error)
	SetACL(path string, acl []zk.ACL, version int32) (*zk.Stat, error)
}
