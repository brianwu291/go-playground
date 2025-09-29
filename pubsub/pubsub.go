package pubsub

import (
	"fmt"
	"sync"
	"time"
)

type PubSubFlash struct{}

const (
	maxSize = 3
)

func NewPubSubFlash() *PubSubFlash {
	return &PubSubFlash{}
}

func (a *PubSubFlash) mockFlashFile(content []string) {
	if len(content) > 0 {
		time.Sleep(time.Second * 3)
		fmt.Printf("file content: %+v\n", content)
	}
}

func (a *PubSubFlash) Producer(content []string) chan string {
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

func (a *PubSubFlash) Consumer(buffers <-chan string) {
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

func Demo() {
	pubsub := NewPubSubFlash()
	content := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	buffers := pubsub.Producer(content)
	contentC := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	buffersC := pubsub.Producer(contentC)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		pubsub.Consumer(buffers)
	}()
	go func() {
		defer wg.Done()
		pubsub.Consumer(buffersC)
	}()
	wg.Wait()
	fmt.Printf("all done!\n")
}
