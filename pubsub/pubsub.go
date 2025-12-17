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

func (a *PubSubFlash) mockFlashFile(content []rune) {
	if len(content) > 0 {
		time.Sleep(time.Second * 3)
		fmt.Printf("file content: %+v\n", string(content))
	}
}

func (a *PubSubFlash) Producer(content []rune) chan rune {
	var wg sync.WaitGroup
	result := make(chan rune, len(content))
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

func (a *PubSubFlash) Consumer(buffers <-chan rune) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	var content []rune

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
	content := []rune("abcdefghij")
	buffers := pubsub.Producer(content)
	contentC := []rune("1234567890")
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
	fmt.Printf("done!\n")
}
