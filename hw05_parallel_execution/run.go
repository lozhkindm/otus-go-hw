package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrWrongWorkersNumber  = errors.New("wrong workers number")
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (result error) {
	if n < 1 {
		return ErrWrongWorkersNumber
	}
	if m < 1 {
		m = 1
	}

	errCount := 0

	wg := &sync.WaitGroup{}
	wg.Add(n)

	taskCh := make(chan Task)
	errCh := make(chan struct{}, m)

	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup, taskCh <-chan Task, errCh chan<- struct{}) {
			for task := range taskCh {
				if err := task(); err != nil {
					errCh <- struct{}{}
				}
			}
			wg.Done()
		}(wg, taskCh, errCh)
	}

	for _, task := range tasks {
		select {
		case <-errCh:
			errCount++
		default:
		}

		if errCount == m {
			result = ErrErrorsLimitExceeded
			break
		}

		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	return result
}
