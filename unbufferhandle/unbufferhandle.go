package unbufferhandle

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type UnBufferHandle struct{}

func NewUnBufferHandle() *UnBufferHandle {
	return &UnBufferHandle{}
}

func (u *UnBufferHandle) Start(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var wg sync.WaitGroup

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
			case <-timeoutCtx.Done():
				fmt.Printf("times up!\n")
				return
			}
		}
	}()

	wg.Wait()
}
