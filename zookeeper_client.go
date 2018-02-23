package curator

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	DefaultSessionTimeout    = 60 * time.Second
	DefaultConnectionTimeout = 15 * time.Second
)

var (
	ErrClientClosed = errors.New("curator: ZookeeperClient had been closed")
)

type ZookeeperClient struct {
	*connectionState
	started           int32
	retryPolicy       RetryPolicy
	quit              chan struct{}
	connectionTimeout time.Duration
}

func NewZookeeperClient(
	factory ZookeeperFactory,
	ensemble EnsembleProvider,
	sessionTimeout, connectTimeout time.Duration,
	retryPolicy RetryPolicy, canBeReadOnly bool) (*ZookeeperClient, error) {

	if sessionTimeout < connectTimeout {
		Log.Warnln("curator: session timeout is less than connection timeout")
	}

	if ensemble == nil {
		return nil, errors.New("curator: ensemble cannot be nil")
	}

	if retryPolicy == nil {
		return nil, errors.New("curator: retryPolicy cannot be nil")
	}

	state := newConnectionState(factory, ensemble, sessionTimeout, connectTimeout, canBeReadOnly)

	client := &ZookeeperClient{
		connectionState:   state,
		retryPolicy:       retryPolicy,
		connectionTimeout: connectTimeout,
		quit:              make(chan struct{}, 1),
	}

	return client, nil
}

func (c *ZookeeperClient) Start() error {
	if !atomic.CompareAndSwapInt32(&c.started, 0, 1) {
		return errors.New("curator: ZookeeperClient already started")
	}
	return c.connectionState.start()
}

func (c *ZookeeperClient) Close() error {
	if atomic.CompareAndSwapInt32(&c.started, 1, 0) {
		c.connectionState.Close()
	}
	return nil
}

func (c *ZookeeperClient) GetConn() Conn {
	if atomic.LoadInt32(&c.started) != 1 {
		return dummyConn{ErrClientClosed}
	}

	conn, err := c.connectionState.getConn()
	if err != nil {
		conn = dummyConn{err}
	}
	return conn
}

func (c *ZookeeperClient) GetRetryPolicy() RetryPolicy {
	return c.retryPolicy
}

func (c *ZookeeperClient) BlockUntilConnectedOrTimeout() {
	if atomic.LoadInt32(&c.started) != 1 {
		return
	}

	timer := time.NewTimer(1 * time.Second)
	waitTime := c.connectionTimeout
	for !c.connectionState.IsConnected() && waitTime > 0 {
		w := make(chan zk.State, 1)
		watcher := NewWatcher(func(event zk.Event) {
			select {
			case w <- event.State:
			default:
			}
		})

		c.AddWatcher(watcher)

		start := time.Now()
		quit := false
		select {
		case <-c.quit:
			quit = true
		case <-w:
		case <-timer.C:
		}
		c.DelWatcher(watcher)

		if quit {
			break
		}

		waitTime -= time.Since(start)
	}
}
