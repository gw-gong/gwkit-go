package common

import (
	"context"
	"fmt"
	"time"
)

const (
	DefaultInterval = 50 * time.Millisecond
	MaxBackoff      = 10 * time.Second
)

// WithRetry is a function that retries a function until it returns false or the context is canceled.
// maxTries is the maximum number of tries, interval is the interval between tries.
// allTries is the number of tries, allErrs is the "returned error" of all tries, lastErr is the error of the last try.
func WithRetry(ctx context.Context, operationName string, f func() (shouldRetry bool, err error), maxTries int, interval time.Duration) (
	allTries int, allErrs []error, lastErr error) {
	if maxTries <= 0 {
		maxTries = 1
	}
	if interval <= 0 {
		interval = DefaultInterval
	}

	allErrs = make([]error, 0, maxTries)
	for i := 0; i < maxTries; {
		shouldRetry, err := f()
		lastErr = err
		if err != nil {
			allErrs = append(allErrs, fmt.Errorf("%s try %d: %w", operationName, i+1, lastErr))
		}
		if !shouldRetry {
			return i + 1, allErrs, lastErr
		}

		// The last loop does not need to wait, just return
		i++
		if i >= maxTries {
			break
		}

		select {
		case <-ctx.Done():
			lastErr = fmt.Errorf("%s try %d: context canceled: %w", operationName, i+1, ctx.Err())
			allErrs = append(allErrs, lastErr)
			return i + 1, allErrs, lastErr
		default:
			backoff := time.Duration(1<<uint(i-1)) * interval
			if backoff > MaxBackoff {
				backoff = MaxBackoff
			}
			time.Sleep(backoff)
		}
	}
	return maxTries, allErrs, lastErr
}
