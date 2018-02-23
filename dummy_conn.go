package curator

import "github.com/samuel/go-zookeeper/zk"

type dummyConn struct {
	err error
}

func (d dummyConn) Close() {
}

func (d dummyConn) Get(path string) (data []byte, stat *zk.Stat, err error) {
	err = d.err
	return
}

func (d dummyConn) GetW(path string) (data []byte, stat *zk.Stat, watch <-chan zk.Event, err error) {
	err = d.err
	return
}

func (d dummyConn) Children(path string) (children []string, stat *zk.Stat, err error) {
	err = d.err
	return
}

func (d dummyConn) ChildrenW(path string) (children []string, stat *zk.Stat, watch <-chan zk.Event, err error) {
	err = d.err
	return
}

func (d dummyConn) Exists(path string) (exist bool, stat *zk.Stat, err error) {
	err = d.err
	return
}

func (d dummyConn) ExistsW(path string) (exist bool, stat *zk.Stat, watch <-chan zk.Event, err error) {
	err = d.err
	return
}

func (d dummyConn) Create(path string, value []byte, flags int32, aclv []zk.ACL) (pathCreated string, err error) {
	err = d.err
	return
}

func (d dummyConn) CreateProtectedEphemeralSequential(path string, value []byte, aclv []zk.ACL) (pathCreated string, err error) {
	err = d.err
	return
}

func (d dummyConn) Set(path string, value []byte, version int32) (stat *zk.Stat, err error) {
	err = d.err
	return
}

func (d dummyConn) Delete(path string, version int32) (err error) {
	err = d.err
	return
}

func (d dummyConn) GetACL(path string) ([]zk.ACL, *zk.Stat, error) {
	return nil, nil, d.err
}

func (d dummyConn) SetACL(path string, aclv []zk.ACL, version int32) (*zk.Stat, error) {
	return nil, d.err
}
