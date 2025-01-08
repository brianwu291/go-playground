package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// "time"

	// ratelimiter "github.com/brianwu291/go-playground/ratelimiter"
	workerpool "github.com/brianwu291/go-playground/workerpool"
)

func main() {
	ctx := context.Background()

	// rateLimiter := ratelimiter.NewRateLimiter(10, time.Second * 60)
	pool := workerpool.NewWorkerPool[int]()
	var jobs []workerpool.Job[int]
	for i := 0; i < 10; i += 1 {
		oneJob := func(timeoutCtx context.Context) (int, error) {
			sleepTime := time.Second * time.Duration(i+30)
			fmt.Printf("sleep: %+v\n", sleepTime)
			time.Sleep(sleepTime)
			select {
			case <-timeoutCtx.Done():
				return 0, fmt.Errorf("job timeout: %+v", i)
			default:
				res := i + 1
				return res, nil
			}
		}
		jobs = append(jobs, oneJob)
	}

	result, err := pool.Process(ctx, jobs)
	fmt.Printf("result %+v\n", result)
	fmt.Printf("err %+v\n", err)

	// firstKey := "key1"
	// for i := 0; i < 11; i += 1 {
	// 	ok := rateLimiter.Access(firstKey)
	// 	// fmt.Printf("ok: %+v on: %d\n", ok, i)
	// }
	// rest := rateLimiter.GetRestTime(firstKey)
	// fmt.Printf("rest time %+v\n", rest)

	portStr := "8999"
	fmt.Printf("listening port %+v\n", portStr)

	http.ListenAndServe(":"+portStr, nil)
}
