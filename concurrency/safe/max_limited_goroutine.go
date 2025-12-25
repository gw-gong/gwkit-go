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

var (
	ErrLimiterClosed = errors.New("limiter is closed")
	ErrNoWaitTimeout = errors.New("max concurrent limit reached and wait timeout is disabled")
)

type AsyncFuncLimiterConfig struct {
	MaxConcurrent int `json:"max_concurrent" yaml:"max_concurrent" mapstructure:"max_concurrent"`
	WaitTimeoutMs int `json:"wait_timeout_ms" yaml:"wait_timeout_ms" mapstructure:"wait_timeout_ms"`
}

const (
	defaultMaxConcurrent = 1000
	defaultWaitTimeoutMs = 500
)

type AsyncFuncLimiter struct {
	semChan     chan struct{}
	waitTimeout time.Duration
	wg          sync.WaitGroup
	closed      atomic.Bool
}

func NewAsyncFuncLimiter(cfg *AsyncFuncLimiterConfig) *AsyncFuncLimiter {
	if cfg == nil {
		cfg = &AsyncFuncLimiterConfig{
			MaxConcurrent: defaultMaxConcurrent,
			WaitTimeoutMs: defaultWaitTimeoutMs,
		}
	}
	if cfg.MaxConcurrent <= 0 {
		cfg.MaxConcurrent = defaultMaxConcurrent
	}
	if cfg.WaitTimeoutMs <= 0 {
		cfg.WaitTimeoutMs = defaultWaitTimeoutMs
	}
	return &AsyncFuncLimiter{
		semChan:     make(chan struct{}, cfg.MaxConcurrent),
		waitTimeout: time.Duration(cfg.WaitTimeoutMs) * time.Millisecond,
		wg:          sync.WaitGroup{},
		closed:      atomic.Bool{},
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
		if afl.waitTimeout == 0 {
			return ErrNoWaitTimeout
		}

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
