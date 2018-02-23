package curator

import "github.com/samuel/go-zookeeper/zk"

func (c *ZookeeperClient) Get(path string) (data []byte, stat *zk.Stat, err error) {
	CallWithRetryLoop(c, func() error {
		data, stat, err = c.GetConn().Get(path)
		return err
	})
	return
}

func (c *ZookeeperClient) GetW(path string) (data []byte, stat *zk.Stat, watch <-chan zk.Event, err error) {
	CallWithRetryLoop(c, func() error {
		data, stat, watch, err = c.GetConn().GetW(path)
		return err
	})
	return
}

func (c *ZookeeperClient) Children(path string) (children []string, stat *zk.Stat, err error) {
	CallWithRetryLoop(c, func() error {
		children, stat, err = c.GetConn().Children(path)
		return err
	})
	return
}

func (c *ZookeeperClient) ChildrenW(path string) (children []string, stat *zk.Stat, watch <-chan zk.Event, err error) {
	CallWithRetryLoop(c, func() error {
		children, stat, watch, err = c.GetConn().ChildrenW(path)
		return err
	})
	return
}

func (c *ZookeeperClient) Exists(path string) (exist bool, stat *zk.Stat, err error) {
	CallWithRetryLoop(c, func() error {
		exist, stat, err = c.GetConn().Exists(path)
		return err
	})
	return
}

func (c *ZookeeperClient) ExistsW(path string) (exist bool, stat *zk.Stat, watch <-chan zk.Event, err error) {
	CallWithRetryLoop(c, func() error {
		exist, stat, watch, err = c.GetConn().ExistsW(path)
		return err
	})
	return
}

func (c *ZookeeperClient) Create(path string, value []byte, flags int32, aclv []zk.ACL) (pathCreated string, err error) {
	CallWithRetryLoop(c, func() error {
		pathCreated, err = c.GetConn().Create(path, value, flags, aclv)
		return err
	})
	return
}

func (c *ZookeeperClient) CreateProtectedEphemeralSequential(path string, value []byte, aclv []zk.ACL) (pathCreated string, err error) {
	CallWithRetryLoop(c, func() error {
		pathCreated, err = c.GetConn().CreateProtectedEphemeralSequential(path, value, aclv)
		return err
	})
	return
}

func (c *ZookeeperClient) Set(path string, value []byte, version int32) (stat *zk.Stat, err error) {
	CallWithRetryLoop(c, func() error {
		stat, err = c.GetConn().Set(path, value, version)
		return err
	})
	return
}

func (c *ZookeeperClient) Delete(path string, version int32) (err error) {
	CallWithRetryLoop(c, func() error {
		err = c.GetConn().Delete(path, version)
		return err
	})
	return
}

func (c *ZookeeperClient) GetACL(path string) (acl []zk.ACL, stat *zk.Stat, err error) {
	CallWithRetryLoop(c, func() error {
		acl, stat, err = c.GetConn().GetACL(path)
		return err
	})
	return
}

func (c *ZookeeperClient) SetACL(path string, acl []zk.ACL, version int32) (stat *zk.Stat, err error) {
	CallWithRetryLoop(c, func() error {
		stat, err = c.GetConn().SetACL(path, acl, version)
		return err
	})
	return
}
