package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrWrongWorkersNumber  = errors.New("wrong routines number")
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var result error

	if n < 1 {
		return ErrWrongWorkersNumber
	}
	if m < 1 {
		return ErrErrorsLimitExceeded
	}

	taskCount := 0
	errCount := 0

	wg := &sync.WaitGroup{}
	wg.Add(n)

	taskCh := make(chan Task)
	errCh := make(chan struct{}, m)
	stopCh := make(chan struct{}, n)

	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup, taskCh <-chan Task, errCh chan<- struct{}, stopCh <-chan struct{}) {
			for {
				select {
				case <-stopCh:
					wg.Done()
					return
				case task := <-taskCh:
					if err := task(); err != nil {
						errCh <- struct{}{}
					}
				}
			}
		}(wg, taskCh, errCh, stopCh)
	}

	for {
		if taskCount == len(tasks) {
			break
		}

		select {
		case <-errCh:
			errCount++
		default:
		}

		if errCount == m {
			result = ErrErrorsLimitExceeded
			break
		}

		taskCh <- tasks[taskCount]
		taskCount++
	}

	for i := 0; i < n; i++ {
		stopCh <- struct{}{}
	}

	wg.Wait()

	return result
}
