package unbufferhandle

import (
	"fmt"
	"sync"
	"time"
)

type UnBufferHandle struct{}

func NewUnBufferHandle() *UnBufferHandle {
	return &UnBufferHandle{}
}

func (u *UnBufferHandle) Start() {
	var wg sync.WaitGroup

	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()

	done := make(chan struct{})
	res := make(chan int)
	wg.Add(1)
	go func() {
		defer close(res)
		defer wg.Done()
		for i := 0; i < 5; i += 1 {
			select {
			case <-done:
				fmt.Printf("stop send at: %+v\n", i)
				return
			default:
				res <- i
				time.Sleep(time.Second * 3)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer close(done)
		defer wg.Done()
		for {
			select {
			case rr, ok := <-res:
				if !ok {
					return
				}
				fmt.Printf("rr: %+v\n", rr)
			case <-timer.C:
				fmt.Printf("times up!\n")
				return
			}
		}
	}()

	wg.Wait()
}
