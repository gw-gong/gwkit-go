package safe

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gw-gong/gwkit-go/util"
)

var ErrLimiterClosed = errors.New("limiter is closed")

type AsyncFuncLimiter struct {
	semChan     chan struct{}
	waitTimeout time.Duration
	wg          sync.WaitGroup
	closed      atomic.Bool
}

func NewAsyncFuncLimiter(maxConcurrent int, waitTimeoutMs int) *AsyncFuncLimiter {
	return &AsyncFuncLimiter{
		semChan:     make(chan struct{}, maxConcurrent),
		waitTimeout: time.Duration(waitTimeoutMs) * time.Millisecond,
	}
}

func (afl *AsyncFuncLimiter) Close() {
	afl.closed.Store(true)
	afl.wg.Wait()
}

func (afl *AsyncFuncLimiter) Async(ctx context.Context, fn func()) error {
	if afl.closed.Load() {
		return ErrLimiterClosed
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case afl.semChan <- struct{}{}:
		if afl.closed.Load() {
			<-afl.semChan
			return ErrLimiterClosed
		}
		afl.executeAsync(fn)
		return nil
	default:
		timeout := time.NewTimer(afl.waitTimeout)
		defer timeout.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return fmt.Errorf("wait timeout(%v)", afl.waitTimeout)
		case afl.semChan <- struct{}{}:
			if afl.closed.Load() {
				<-afl.semChan
				return ErrLimiterClosed
			}
			afl.executeAsync(fn)
			return nil
		}
	}
}

func (afl *AsyncFuncLimiter) executeAsync(fn func()) {
	afl.wg.Add(1)
	go util.WithRecover(func() {
		defer func() {
			<-afl.semChan
			afl.wg.Done()
		}()
		fn()
	})
}
