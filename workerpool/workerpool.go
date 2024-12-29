package workerpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WorkerPool[T any] struct {
	Config WorkerPoolConfig
}

type WorkerPoolConfig struct {
	Size    int
	Timeout time.Duration
}

type Job[T any] func() (T, error)

var defaultConfig = WorkerPoolConfig{
	Size:    5,
	Timeout: time.Second * 30,
}

func NewWorkerPool[T any](options ...func(*WorkerPoolConfig)) *WorkerPool[T] {
	config := defaultConfig
	for _, option := range options {
		option(&config)
	}
	return &WorkerPool[T]{Config: config}
}

func (w *WorkerPool[T]) Process(ctx context.Context, jobs []Job[T]) ([]T, []error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, w.Config.Timeout)
	defer cancel()

	howManyJobs := len(jobs)
	internalJobs := make(chan Job[T], howManyJobs)
	results := make(chan T, howManyJobs)
	errors := make(chan error, howManyJobs)

	for i := 0; i < howManyJobs; i += 1 {
		internalJobs <- jobs[i]
	}
	close(internalJobs)

	var wg sync.WaitGroup
	for i := 0; i < w.Config.Size; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iJob := range internalJobs {
				select {
				case <-timeoutCtx.Done():
					errors <- fmt.Errorf("worker %d: timeout", i)
					return
				default:
					res, err := iJob()
					fmt.Printf("processed job: %+v with goroutine: %+v \n", res, i+1)
					results <- res
					if err != nil {
						errors <- err
					}
				}
			}
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	var res []T
	var errs []error

	for result := range results {
		res = append(res, result)
	}

	for err := range errors {
		errs = append(errs, err)
	}

	return res, errs
}
