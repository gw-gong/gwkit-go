package pool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gw-gong/gwkit-go/util/common"
)

const (
	DefaultTimeoutSubmit = time.Second * 5
)

type WorkerPool interface {
	Submit(work Work) error
	QueueLength() int
	Close()
}

type workerPoolImpl struct {
	wg            sync.WaitGroup
	closed        int32
	workChan      chan Work
	timeoutSubmit time.Duration
}

type option func(wp *workerPoolImpl)

func WithTimeoutSubmit(timeoutSubmit time.Duration) option {
	return func(wp *workerPoolImpl) {
		if timeoutSubmit > 0 {
			wp.timeoutSubmit = timeoutSubmit
		}
	}
}

func NewWorkerPool(channelSize, workerPoolSize int, opts ...option) (WorkerPool, error) {
	if channelSize <= 0 {
		return nil, errors.New("channelSize must be greater than 0")
	}
	if workerPoolSize <= 0 {
		return nil, errors.New("workerPoolSize must be greater than 0")
	}

	wp := &workerPoolImpl{
		wg:            sync.WaitGroup{},
		closed:        0,
		workChan:      make(chan Work, channelSize),
		timeoutSubmit: DefaultTimeoutSubmit,
	}

	for _, opt := range opts {
		opt(wp)
	}

	wp.run(workerPoolSize)

	return wp, nil
}

func (wp *workerPoolImpl) run(workerPoolSize int) {
	for i := 0; i < workerPoolSize; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for work := range wp.workChan {
				common.WithRecover(work.Do)
			}
		}()
	}
}

func (wp *workerPoolImpl) Submit(work Work) error {
	if atomic.LoadInt32(&wp.closed) == 1 {
		return errors.New("worker pool is closed")
	}

	timeout := time.NewTimer(wp.timeoutSubmit)
	defer timeout.Stop()

	select {
	case wp.workChan <- work:
	case <-timeout.C:
		return fmt.Errorf("timeout(%v) to submit work, workChan is full, queue length: %d", wp.timeoutSubmit, wp.QueueLength())
	}
	return nil
}

func (wp *workerPoolImpl) QueueLength() int {
	return len(wp.workChan)
}

func (wp *workerPoolImpl) Close() {
	if atomic.CompareAndSwapInt32(&wp.closed, 0, 1) {
		close(wp.workChan)
		wp.wg.Wait()
	}
}
