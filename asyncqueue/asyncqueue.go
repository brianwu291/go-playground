package asyncqueue

import (
	"context"
	"fmt"
	"sync"
)

type Task func() error

type AsyncQueue struct {
	tasks chan Task
	once  sync.Once
}

func New(concurrency uint) (*AsyncQueue, error) {
	if concurrency == 0 {
		return nil, fmt.Errorf("asyncqueue: concurrency must be greater than zero")
	}
	q := &AsyncQueue{
		tasks: make(chan Task, concurrency*10),
	}
	for range concurrency {
		go func() {
			for t := range q.tasks {
				t()
			}
		}()
	}
	return q, nil
}

func (q *AsyncQueue) Add(ctx context.Context, task Task) error {
	select {
	case q.tasks <- task:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("asyncqueue: %w", ctx.Err())
	}
}

func (q *AsyncQueue) Close() {
	q.once.Do(func() {
		close(q.tasks)
	})
}
