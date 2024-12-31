package main

import (
	"fmt"
	"net/http"
	"sync"

	// "time"

	// workerpool "github.com/brianwu291/go-playground/workerpool"
	// unbufferhandle "github.com/brianwu291/go-playground/unbufferhandle"
	adder "github.com/brianwu291/go-playground/adder"
)

func main() {
	// ctx := context.Background()
	// unBufferHandle := unbufferhandle.NewUnBufferHandle()
	// unBufferHandle.Start(ctx)

	var wg sync.WaitGroup
	adder := adder.NewAdder(&adder.AdderConfig{Size: 5})
	for i := 0; i < 1000; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			adder.Inc(1)
		}()
	}
	first := adder.GetCurrentValue()
	fmt.Printf("first: %d \n", first)
	wg.Wait()
	second := adder.GetCurrentValue()
	fmt.Printf("second: %d \n", second)
	adder.Close()

	portStr := "8999"
	fmt.Printf("listening port %+v\n", portStr)

	http.ListenAndServe(":"+portStr, nil)
}
