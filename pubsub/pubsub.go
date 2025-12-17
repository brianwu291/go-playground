package pubsub

import (
	"fmt"
	"sync"
	"time"
)

type PubSubFlash struct {
	content        chan byte
	started        bool
	closed         bool
	mu             sync.RWMutex
	flushFunc      func(content []byte)
	maxContentSize int
	interval       time.Duration
	produceTimeout time.Duration
}

const (
	MaxSize = 3
)

var (
	ErrNotStartedYet            = fmt.Errorf("pub sub not started yet")
	ErrProduceTimeout           = fmt.Errorf("produce content timeout")
	ErrProduceFailedDueToClosed = fmt.Errorf("produce content not in time")
)

func NewPubSubFlash(maxSize int, interval, produceTimeout time.Duration, flush func(content []byte)) *PubSubFlash {
	c := make(chan byte, maxSize)
	return &PubSubFlash{
		content:        c,
		flushFunc:      flush,
		maxContentSize: maxSize,
		interval:       interval,
		produceTimeout: produceTimeout,
	}
}

func MockFlushFile(content []byte) {
	if len(content) > 0 {
		time.Sleep(time.Second * 3)
		fmt.Printf("file content: %+v\n", string(content))
	}
}

func (a *PubSubFlash) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.closed {
		a.closed = true
		close(a.content)
	}
}

func (a *PubSubFlash) isClosed() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.closed
}

func (a *PubSubFlash) isStarted() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.started
}

func (a *PubSubFlash) Produce(content []byte) error {
	if !a.isStarted() {
		return ErrNotStartedYet
	}

	timer := time.NewTimer(a.produceTimeout)
	defer timer.Stop()

	for _, cb := range content {
		if a.isClosed() {
			return ErrProduceFailedDueToClosed
		}
		select {
		case a.content <- cb:
		case <-timer.C:
			return ErrProduceTimeout
		}
	}
	return nil
}

func (a *PubSubFlash) consume() {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	var content []byte

	for {
		select {
		case buf, ok := <-a.content:
			if !ok { // closed
				a.flushFunc(content)
				return
			}
			content = append(content, buf)
			if len(content) >= a.maxContentSize {
				a.flushFunc(content)
				content = content[:0]
			}
		case <-ticker.C:
			a.flushFunc(content)
			content = content[:0]
		}
	}
}

func (a *PubSubFlash) Start() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.started {
		a.started = true
		go a.consume()
	}
}

func Demo() {
	max := 6
	interval := time.Second * 3
	produceTimeout := time.Second * 10
	pubsub := NewPubSubFlash(max, interval, produceTimeout, MockFlushFile)
	pubsub.Start()
	content := []byte("abcdefghij")
	err := pubsub.Produce(content)
	if err != nil {
		fmt.Printf("err on p1: %+v", err)
	}
	contentC := []byte("1234567890")
	err = pubsub.Produce(contentC)
	if err != nil {
		fmt.Printf("err on p2: %+v", err)
	}

	time.Sleep(10 * time.Minute)
	fmt.Println("demo pubsub finished")
}
