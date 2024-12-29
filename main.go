package main

import (
	// "context"
	"fmt"
	"net/http"

	// workerpool "github.com/brianwu291/go-playground/workerpool"
	unbufferhandle "github.com/brianwu291/go-playground/unbufferhandle"
)

func main() {
	// ctx := context.Background()
	// pool := workerpool.NewWorkerPool[int](func(c *workerpool.WorkerPoolConfig) {
	//   c.Timeout = time.Second * 15
	// })
	// var jobs []workerpool.Job[int]
	// for i := 0; i < 20; i += 1 {
	// 	jobNum := i + 1
	// 	var job workerpool.Job[int] = func() (int, error) {
	// 		time.Sleep(time.Second * 4)
	// 		return jobNum, nil
	// 	}
	// 	jobs = append(jobs, job)
	// }
	// results, errors := pool.Process(ctx, jobs)
	// fmt.Printf("results: %+v \n", results)
	// fmt.Printf("errors: %+v \n", errors)

	unBufferHandle := unbufferhandle.NewUnBufferHandle()
	unBufferHandle.Start()

	portStr := "8999"
	fmt.Printf("listening port %+v\n", portStr)

	http.ListenAndServe(":"+portStr, nil)
}
