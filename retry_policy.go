package curator

import (
	"math/rand"
	"time"
)

type RetrySleeper interface {
	Sleep(duration time.Duration) error
}
type RetrySleeperFunc func(duration time.Duration) error

func (r RetrySleeperFunc) Sleep(duration time.Duration) error {
	return r(duration)
}

func NopRetrySleeper() RetrySleeper {
	return RetrySleeperFunc(func(time.Duration) error { return nil })
}

type RetryPolicy interface {
	AllowRetry(count int, elapsed time.Duration, sleeper RetrySleeper) bool
}
type RetryPolicyFunc func(count int, elapsed time.Duration, sleeper RetrySleeper) bool

func (r RetryPolicyFunc) AllowRetry(count int, elapsed time.Duration, sleeper RetrySleeper) bool {
	return r(count, elapsed, sleeper)
}

type RetryForever struct {
	Duration time.Duration
}

func NewRetryForever(duration time.Duration) RetryForever {
	return RetryForever{duration}
}

func (r RetryForever) AllowRetry(count int, elapsed time.Duration, sleeper RetrySleeper) bool {
	return sleeper.Sleep(r.Duration) == nil
}

type RetryNTimes struct {
	N        int
	Duration time.Duration
}

func NewRetryNTimes(n int, duration time.Duration) RetryNTimes {
	return RetryNTimes{n, duration}
}

func (r RetryNTimes) AllowRetry(count int, elapsed time.Duration, sleeper RetrySleeper) bool {
	allowed := false
	if count < r.N {
		allowed = sleeper.Sleep(r.Duration) == nil
	}
	return allowed
}

type ExponentialBackoffRetry struct {
	N            int
	BaseDuration time.Duration
	MaxDuration  time.Duration
}

func NewExponentialBackoffRetry(n int, baseDuration, maxDuration time.Duration) ExponentialBackoffRetry {
	return ExponentialBackoffRetry{
		N:            n,
		BaseDuration: baseDuration,
		MaxDuration:  maxDuration}
}

func (r ExponentialBackoffRetry) AllowRetry(count int, elapsed time.Duration, sleeper RetrySleeper) bool {
	maxRetries := 29
	if maxRetries < r.N {
		maxRetries = r.N
	}

	allowed := false
	if count < r.N {
		n := rand.Intn(int(1 << uint(count+1)))
		if n < 1 {
			n = 1
		}
		duration := r.BaseDuration * time.Duration(n)
		if duration > r.MaxDuration {
			duration = r.MaxDuration
		}

		allowed = sleeper.Sleep(duration) == nil
	}
	return allowed
}
