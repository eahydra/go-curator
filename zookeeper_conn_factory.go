package curator

import (
	"net"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZookeeperFactory interface {
	Create(connectString string, sessionTimeout, connectTimeout time.Duration, readOnly bool) (*zk.Conn, <-chan zk.Event, error)
}

type defaultZookeeperFactory struct{}

var DefaultZookeeperFactory = defaultZookeeperFactory{}

func (defaultZookeeperFactory) Create(connectString string, sessionTimeout, connectTimeout time.Duration, readOnly bool) (*zk.Conn, <-chan zk.Event, error) {
	return zk.Connect(strings.Split(connectString, ","), sessionTimeout, zk.WithDialer(dialWithTimeout(connectTimeout)))
}

func dialWithTimeout(timeout time.Duration) zk.Dialer {
	return func(network, address string, _ time.Duration) (net.Conn, error) {
		return net.DialTimeout(network, address, timeout)
	}
}
