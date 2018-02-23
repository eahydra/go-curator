package curator

import (
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type connHolder struct {
	conn              Conn
	factory           ZookeeperFactory
	ensemble          EnsembleProvider
	canBeReadOnly     bool
	sessionTimeout    time.Duration
	connectionTimeout time.Duration
	connectionString  string
	mutex             sync.Mutex
	processEvent      func(zk.Event)
}

func (c *connHolder) getConn() Conn {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.conn
}

func (c *connHolder) hasNewConnectionString() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.connectionString != c.ensemble.GetConnectionString()
}

func (c *connHolder) closeAndClear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *connHolder) closeAndReset() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}

	c.connectionString = c.ensemble.GetConnectionString()
	conn, watch, err := c.factory.Create(c.connectionString, c.sessionTimeout, c.connectionTimeout, c.canBeReadOnly)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.watch(watch)

	return nil
}

func (c *connHolder) watch(eventChan <-chan zk.Event) {
	for {
		select {
		case event, ok := <-eventChan:
			Log.Debugf("curator: connHolder got event:%+v", event)
			if !ok {
				return
			}
			c.processEvent(event)
		}
	}
}
