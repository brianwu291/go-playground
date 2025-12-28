package pb

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// chan d with init buff (should > 0), default 50
// store ctx as inherit outside ctx done
// Start to run inner consume
// Produce response for send d into chan
// consume keep extract d from chan, at the same time listen ctx, timeout, ...etc. Graceful shut down to call flush
// pass flush method to do "things" u want from those producer

const (
	DefaultBufSize = 50
)

var (
	ErrPBNotStarted = fmt.Errorf("PB not started yet")
)

type PB struct {
	started      bool
	mu           sync.RWMutex
	d            chan byte
	bufferSize   int
	flush        func(d []byte)
	flushFreq    time.Duration
	flushSize    int
	flushTimeout time.Duration
}

func MockFlushFile(content []byte) {
	if len(content) > 0 {
		time.Sleep(time.Second * 3)
		fmt.Printf("file content: %+v\n", string(content))
	}
}

func NewPB(bufSize int, flush func(d []byte), flushFreq, flushTimeout time.Duration, flushSize int) *PB {
	if bufSize <= 1 {
		bufSize = DefaultBufSize
	}
	if flushSize >= bufSize {
		flushSize = bufSize - 1
	}
	initChan := make(chan byte, bufSize)
	return &PB{
		d:            initChan,
		flush:        flush,
		flushFreq:    flushFreq,
		flushSize:    flushSize,
		flushTimeout: flushTimeout,
	}
}

func (pb *PB) Start(ctx context.Context) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	if pb.started {
		return
	}
	pb.started = true
	go pb.consume(ctx)
}

func (pb *PB) Stop() {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	if pb.started {
		pb.started = false
		close(pb.d)
	}
}

func (pb *PB) isStarted() bool {
	pb.mu.RLock()
	defer pb.mu.RUnlock()
	return pb.started
}

func (pb *PB) Produce(data []byte) error {
	for _, d := range data {
		if !pb.isStarted() {
			return ErrPBNotStarted
		}
		pb.d <- d // maybe add a timeout config and leverage select here to handle
	}
	return nil
}

func (pb *PB) consume(ctx context.Context) {
	ticker := time.NewTicker(pb.flushFreq)
	defer ticker.Stop()

	content := make([]byte, 0, pb.bufferSize)

	for {
		select {
		case data, ok := <-pb.d:
			if !ok {
				pb.executeFlushWithTimeout(content, pb.flushTimeout)
				return
			}
			content = append(content, data)
			if len(content) >= pb.flushSize {
				pb.executeFlushWithTimeout(content, pb.flushTimeout)
				content = content[:0]
			}
		case <-ticker.C:
			pb.executeFlushWithTimeout(content, pb.flushTimeout)
			content = content[:0]
		case <-ctx.Done():
			fmt.Printf("stopped!\n")
			pb.executeFlushWithTimeout(content, pb.flushTimeout)
			pb.Stop()
			return
		}
	}
}

func (pb *PB) executeFlushWithTimeout(d []byte, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		pb.flush(d)
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		fmt.Printf("flush timeout after %v, dropping %d bytes\n", timeout, len(d))
		return
	}
}

func Demo() {
	bufSize := 10
	flushSize := 4
	interval := time.Second * 10
	flushTimeout := time.Second * 5
	pb := NewPB(bufSize, MockFlushFile, interval, flushTimeout, flushSize)
	ctx, cancel := context.WithCancel(context.Background())
	pb.Start(ctx)
	content := []byte("abcdefghij")
	err := pb.Produce(content)
	if err != nil {
		fmt.Printf("err on p1: %+v", err)
	}
	contentC := []byte("1234567890")
	err = pb.Produce(contentC)
	if err != nil {
		fmt.Printf("err on p2: %+v", err)
	}

	time.Sleep(1 * time.Minute)
	cancel()
	fmt.Println("demo pb finished")
	pb.Stop()
}
