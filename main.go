package main

import (
	"fmt"
	"net/http"
	"time"

	ratelimiter "github.com/brianwu291/go-playground/ratelimiter"
)

func main() {
	// ctx := context.Background()

	rateLimiter := ratelimiter.NewRateLimiter(10, time.Second * 60)

	firstKey := "key1"
	for i := 0; i < 11; i += 1 {
		ok := rateLimiter.Access(firstKey)
		fmt.Printf("ok: %+v on: %d\n", ok, i)
	}
	rest := rateLimiter.GetRestTime(firstKey)
	fmt.Printf("rest time %+v\n", rest)

	portStr := "8999"
	fmt.Printf("listening port %+v\n", portStr)

	http.ListenAndServe(":"+portStr, nil)
}
