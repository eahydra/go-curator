package curator

import (
	"errors"
	"testing"
	"time"
)

func TestExponentialBackoffRetry_AllowRetry(t *testing.T) {
	baseDuration := 100 * time.Millisecond
	maxDuration := 1 * time.Second
	retry := NewExponentialBackoffRetry(30, baseDuration, maxDuration)
	for i := 0; i < 3; i++ {
		if !retry.AllowRetry(i, 0, RetrySleeperFunc(func(duration time.Duration) error {
			if duration < baseDuration || duration > maxDuration {
				return errors.New("invalid duration:" + duration.String())
			}
			time.Sleep(duration)
			return nil
		})) {
			t.Fatal("unexpected AllowRetry")
		}
	}
	if retry.AllowRetry(31, 0, NopRetrySleeper()) {
		t.Fatal("unexpected AllowRetry")
	}
}

func TestRetryNTimes_AllowRetry(t *testing.T) {
	retry := NewRetryNTimes(1, 100*time.Millisecond)
	if !retry.AllowRetry(0, 0, NopRetrySleeper()) {
		t.Fatal("unexpected AllowRetry")
	}
	if retry.AllowRetry(1, 0, NopRetrySleeper()) {
		t.Fatal("unexpected AllowRetry")
	}
	retry = NewRetryNTimes(0, 100*time.Millisecond)
	if retry.AllowRetry(0, 0, NopRetrySleeper()) {
		t.Fatal("unexpected AllowRetry")
	}
}

func TestRetryForever_AllowRetry(t *testing.T) {
	retry := NewRetryForever(100 * time.Millisecond)
	if !retry.AllowRetry(1, 0, NopRetrySleeper()) {
		t.Fatal("unexpected AllowRetry")
	}

	if !retry.AllowRetry(0, 0, NopRetrySleeper()) {
		t.Fatal("unexpected AllowRetry")
	}
}
