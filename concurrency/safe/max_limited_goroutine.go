package safe

import (
	"fmt"
	"sync"
)

const DefaultMaxGoroutine = 1000

// Use a struct to wrap the function as a key, to avoid using the function directly as a map key
type funcKey struct {
	f interface{} // Store the function, use interface{} to avoid direct function comparison
}

var (
	funcGoroutineSem = sync.Map{}
	semMu            sync.Mutex // Protect the creation of the semaphore
)

// Async executes the function f asynchronously, limiting the number of concurrent goroutines for the same function to no more than max.
// If the limit is exceeded, it returns an error.
// !!! Suitable for unimportant asynchronous tasks, because it does not guarantee that the task is successfully submitted and executed !!!
func AsyncUnreliably(f func(), maxGoroutine ...int) error {
	key := funcKey{f: f}
	semInterface, ok := funcGoroutineSem.Load(key)

	// If the semaphore does not exist, create it
	if !ok {
		semMu.Lock()
		defer semMu.Unlock()
		// Double check to avoid concurrent creation
		if semInterface, ok = funcGoroutineSem.Load(key); !ok {
			max := DefaultMaxGoroutine
			if len(maxGoroutine) > 0 && maxGoroutine[0] > 0 {
				max = maxGoroutine[0]
			}
			semInterface = make(chan struct{}, max)
			funcGoroutineSem.Store(key, semInterface)
		}
	}

	sem, ok := semInterface.(chan struct{})
	if !ok {
		return fmt.Errorf("sem is not a channel")
	}

	// Try to acquire the semaphore
	select {
	case sem <- struct{}{}:
	default:
		return fmt.Errorf("exceed max goroutine limit")
	}

	// Start a goroutine to execute the task
	go func() {
		defer func() {
			// Release the semaphore after the task is done
			<-sem
		}()
		f()
	}()

	return nil
}
