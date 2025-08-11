package common

import (
	"context"
	"time"
)

const (
	DefaultInterval = 50 * time.Millisecond
	MaxBackoff      = 10 * time.Second
)

func WithRetry(ctx context.Context, f func() (shouldRetry bool, err error), maxTries int, interval time.Duration) (allTries int, allErrors []error, lastError error) {
	if maxTries <= 0 {
		maxTries = 1
	}
	if interval <= 0 {
		interval = DefaultInterval
	}

	allErrors = make([]error, 0, maxTries)
	for i := 0; i < maxTries; i++ {
		shouldRetry, err := f()
		allErrors = append(allErrors, err)
		if !shouldRetry {
			return i + 1, allErrors, err
		}

		select {
		case <-ctx.Done():
			return i + 1, allErrors, ctx.Err()
		default:
			backoff := time.Duration(1<<uint(i)) * interval
			if backoff > MaxBackoff {
				backoff = MaxBackoff
			}
			time.Sleep(backoff)
		}
	}
	return maxTries, allErrors, allErrors[len(allErrors)-1] // Logically, this will not be out of bounds
}
