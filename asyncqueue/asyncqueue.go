package asyncqueue

import "fmt"

type AsyncQueue struct {
	max         uint
	pushedTasks chan Task
}

type Task func() error

func NewAsyncQueue(limit uint) (*AsyncQueue, error) {
	if limit == 0 {
		return nil, fmt.Errorf("limit should be positive!\n")
	}
	pt := make(chan Task, limit*10)
	aq := &AsyncQueue{
		max:         limit,
		pushedTasks: pt,
	}
	for i := 0; i < int(aq.max); i++ {
		go func() {
			for t := range aq.pushedTasks {
				t()
			}
		}()
	}
	return aq, nil
}

func (aq *AsyncQueue) Add(task Task) {
	aq.pushedTasks <- task
}
