package safe

import "errors"

type AsyncFuncLimiter struct {
	semaphore chan struct{}
}

func NewAsyncFuncLimiter(maxConcurrent int) *AsyncFuncLimiter {
	return &AsyncFuncLimiter{
		semaphore: make(chan struct{}, maxConcurrent),
	}
}

func (fl *AsyncFuncLimiter) Async(fn func()) error {
	select {
	case fl.semaphore <- struct{}{}:
		go func() {
			defer func() { <-fl.semaphore }()
			fn()
		}()
		return nil
	default:
		return errors.New("max concurrent limit exceeded")
	}
}
