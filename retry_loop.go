package curator

import (
	"errors"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type defaultSleeper struct {
	client *ZookeeperClient
}

func (s defaultSleeper) Sleep(duration time.Duration) error {
	select {
	case <-s.client.quit:
		return errors.New("curator: ZookeeperClient had been closed")
	case <-time.After(duration):
		return nil
	}
}

func ShouldRetry(err error) bool {
	switch err {
	default:
		return false
	case zk.ErrSessionExpired, zk.ErrSessionMoved, zk.ErrConnectionClosed, zk.ErrNoServer, zk.ErrClosing, ErrConnectionLoss:
		return true
	}
}

func CallWithRetryLoop(client *ZookeeperClient, operate func() error) (err error) {
	var (
		startTime = time.Now()
		policy    = client.GetRetryPolicy()
		sleeper   = defaultSleeper{client}
	)

	for count := 0; ; count++ {
		client.BlockUntilConnectedOrTimeout()

		err = operate()
		if !ShouldRetry(err) {
			break
		}

		if !policy.AllowRetry(count, time.Since(startTime), sleeper) {
			break
		}
	}
	return
}
