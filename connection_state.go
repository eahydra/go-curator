package curator

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var ErrConnectionLoss = errors.New("curator: connection loss")

const (
	ConnDisconnected = 0
	ConnConnected    = 1
)

type connectionState struct {
	*watcherManager
	holder              *connHolder
	ensemble            EnsembleProvider
	sessionTimeout      time.Duration
	connectionTimeout   time.Duration
	connectionStartTime time.Duration
	connected           int32
	errQueue            *errorQueue
	connStateWatches    []chan<- zk.State
	checkMutex          sync.Mutex
}

func newConnectionState(
	factory ZookeeperFactory,
	ensemble EnsembleProvider,
	sessionTimeout, connectionTimeout time.Duration,
	canBeReadOnly bool) *connectionState {

	state := &connectionState{
		ensemble:          ensemble,
		sessionTimeout:    sessionTimeout,
		connectionTimeout: connectionTimeout,
		errQueue:          newErrorQueue(10),
		watcherManager:    newWatcherManager(),
	}

	holder := &connHolder{
		factory:           factory,
		ensemble:          ensemble,
		canBeReadOnly:     canBeReadOnly,
		sessionTimeout:    sessionTimeout,
		connectionTimeout: connectionTimeout,
		processEvent:      state.processEvent,
	}
	state.holder = holder

	return state
}

func (c *connectionState) start() error {
	if err := c.ensemble.Start(); err != nil {
		return err
	}
	return c.reset()
}

func (c *connectionState) Close() error {
	c.ensemble.Close()
	c.holder.closeAndClear()
	atomic.StoreInt32(&c.connected, ConnDisconnected)
	return nil
}

func (c *connectionState) IsConnected() bool {
	return atomic.LoadInt32(&c.connected) == ConnConnected
}

func (c *connectionState) getConn() (Conn, error) {
	if err := c.errQueue.Pop(); err != nil {
		return nil, err
	}

	if !c.IsConnected() {
		c.checkTimeout()
	}
	return c.holder.getConn(), nil
}

func (c *connectionState) checkTimeout() error {
	c.checkMutex.Lock()
	defer c.checkMutex.Unlock()

	minDuration := c.sessionTimeout
	maxDuration := c.connectionTimeout
	if minDuration > maxDuration {
		minDuration, maxDuration = maxDuration, minDuration
	}

	var err error
	duration := time.Duration(time.Now().UnixNano() - int64(c.connectionStartTime))
	if duration >= minDuration {
		if c.holder.hasNewConnectionString() {
			err = c.reset()

		} else {
			if duration > maxDuration {
				Log.Warnf("curator: connection attemp unsuccessful after %s (greater than max timeout of %s). "+
					"Resetting connection and try again with a newConnection", duration.String(), maxDuration.String())
				err = c.reset()
			} else {
				err = ErrConnectionLoss
			}
		}

	}
	return err
}

func (c *connectionState) reset() error {
	atomic.StoreInt32(&c.connected, ConnDisconnected)
	atomic.StoreInt64((*int64)(&c.connectionStartTime), time.Now().UnixNano())
	return c.holder.closeAndReset()
}

func (c *connectionState) processEvent(event zk.Event) {
	Log.Infof("curator: got conn event: %+v", event)

	checkConnectionString := true

	if event.State == zk.StateConnected {

		atomic.StoreInt32(&c.connected, ConnConnected)
		atomic.StoreInt64((*int64)(&c.connectionStartTime), time.Now().UnixNano())

	} else if event.State == zk.StateExpired {
		if err := c.reset(); err != nil {
			c.errQueue.Push(err)
		}
		checkConnectionString = false
	}

	if checkConnectionString && c.holder.hasNewConnectionString() {
		if err := c.reset(); err != nil {
			c.errQueue.Push(err)
		}
	}

	c.Fire(event)
}
