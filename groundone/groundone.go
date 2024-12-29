package groundone

import (
	"fmt"
	"time"
)

type GroundOne struct{}

const (
	maxSize = 3
)

func NewGroundOne() *GroundOne {
	return &GroundOne{}
}

func (a *GroundOne) mockFlashFile(content []string) {
	fmt.Printf("file content: %+v\n", content)
}

func (a *GroundOne) Producer(content []string) chan string {
	result := make(chan string, len(content))
	go func() {
		for _, contentStr := range content {
				time.Sleep(time.Second * 1)
				result <- contentStr
		}
		close(result)
	}()
	return result
}

func (a *GroundOne) Consumer(buffers <-chan string) {
	var fileContent []string
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case buffer, ok := <-buffers:
			if !ok {
				fmt.Printf("buffer closed!\n")
				if len(fileContent) > 0 {
					a.mockFlashFile(fileContent)
				}
				return
			}
			if len(fileContent) == maxSize {
				fmt.Printf("maxSize reached: ")
				a.mockFlashFile(fileContent)
				fileContent = fileContent[:0]
				fileContent = append(fileContent, buffer)
			} else {
				fileContent = append(fileContent, buffer)
			}
		case <-ticker.C:
			fmt.Printf("ticker ta: ")
			if len(fileContent) > 0 {
				a.mockFlashFile(fileContent)
				fileContent = fileContent[:0]
			}
		}
	}
}
