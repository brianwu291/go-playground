package groundone

import (
	"fmt"
	"sync"
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
	if len(content) > 0 {
		time.Sleep(time.Second * 3)
		fmt.Printf("file content: %+v\n", content)
	}
}

func (a *GroundOne) Producer(content []string) chan string {
	var wg sync.WaitGroup
	result := make(chan string, len(content))
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, contentStr := range content {
			result <- contentStr
		}
		close(result)
	}()
	wg.Wait()
	return result
}

func (a *GroundOne) Consumer(buffers <-chan string) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	var content []string

	for {
		select {
		case buf, ok := <-buffers:
			if !ok { // closed
				a.mockFlashFile(content)
				return
			}
			if len(content) == maxSize {
				a.mockFlashFile(content)
				content = content[:0]
				content = append(content, buf)
				continue
			}
			content = append(content, buf)
		case <-ticker.C:
			fmt.Printf("flash anyway!\n")
			a.mockFlashFile(content)
			content = content[:0]
		}
	}
}
