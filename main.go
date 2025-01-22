package main

import (
	"fmt"
	"net/http"
	"time"

	// ratelimiter "github.com/brianwu291/go-playground/ratelimiter"
	// interview "github.com/brianwu291/go-playground/interview"
	// groundone "github.com/brianwu291/go-playground/groundone"
	realtimechat "github.com/brianwu291/go-playground/realtimechat"
)

func main() {
	// ctx := context.Background()
	// gro := groundone.NewGroundOne()
	// var mockFile []string
	// for i := 0; i < 20; i += 1 {
	// 	content := fmt.Sprintf("content: %+v", i)
	// 	mockFile = append(mockFile, content)
	// }
	// buffers := gro.Producer(mockFile)
	// gro.Consumer(buffers)

	chat := realtimechat.NewRealTimeChat(10)
  chat.Run()
  defer chat.Stop()

	clientOne, _ := chat.AddClient("User1")
	clientTwo, _ := chat.AddClient("User2")
	clientThree, _ := chat.AddClient("User3")

	chat.SendMessage("Hi! 11", clientOne.Id)
	chat.SendMessage("Hi! 22", clientTwo.Id)
	chat.SendMessage("Hi! 33", clientThree.Id)

	
	time.Sleep(time.Second * 10)
	chat.Stop()

	// rateLimiter := ratelimiter.NewRateLimiter(10, time.Second * 60)
	// pool := workerpool.NewWorkerPool[int]()
	// var jobs []workerpool.Job[int]
	// for i := 0; i < 10; i += 1 {
	// 	oneJob := func(timeoutCtx context.Context) (int, error) {
	// 		sleepTime := time.Second * time.Duration(i+30)
	// 		fmt.Printf("sleep: %+v\n", sleepTime)
	// 		time.Sleep(sleepTime)
	// 		select {
	// 		case <-timeoutCtx.Done():
	// 			return 0, fmt.Errorf("job timeout: %+v", i)
	// 		default:
	// 			res := i + 1
	// 			return res, nil
	// 		}
	// 	}
	// 	jobs = append(jobs, oneJob)
	// }

	// result, err := pool.Process(ctx, jobs)
	// fmt.Printf("result %+v\n", result)
	// fmt.Printf("err %+v\n", err)

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
