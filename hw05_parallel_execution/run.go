package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrErrorsLimitExceeded ошибка при достижении лимита ошибок.
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// ErrZeroWorkers ошибка при нулевом количестве воркеров.
var ErrZeroWorkers = errors.New("zero workers")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return ErrZeroWorkers
	}

	ch := make(chan Task, n)
	wg := sync.WaitGroup{}

	errCount := int32(0)
	maxErrors := int32(m)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go run(ch, &wg, &errCount, maxErrors)
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errCount) >= maxErrors {
			break
		}
		ch <- task
	}
	close(ch)

	wg.Wait()

	if atomic.LoadInt32(&errCount) >= maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// run функция-форкер, выполняющая пришедшие в канале задачи.
func run(task chan Task, wg *sync.WaitGroup, errCount *int32, maxErrors int32) {
	defer wg.Done()
	for {
		if t, ok := <-task; ok && atomic.LoadInt32(errCount) < maxErrors {
			if err := t(); err != nil {
				atomic.AddInt32(errCount, 1)
			}
		} else {
			return
		}
	}
}
