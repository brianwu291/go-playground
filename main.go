package main

import (
	"context"
	"fmt"
	"net/http"

	// workerpool "github.com/brianwu291/go-playground/workerpool"
	unbufferhandle "github.com/brianwu291/go-playground/unbufferhandle"
)

func main() {
	ctx := context.Background()
	unBufferHandle := unbufferhandle.NewUnBufferHandle()
	unBufferHandle.Start(ctx)

	portStr := "8999"
	fmt.Printf("listening port %+v\n", portStr)

	http.ListenAndServe(":"+portStr, nil)
}
