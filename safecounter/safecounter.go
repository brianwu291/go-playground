package safecounter

import (
	"sync"
)

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

type RWCounter struct {
	mu    sync.RWMutex
	count int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		mu:    sync.Mutex{},
		count: 0,
	}
}

func NewRWCounter() *RWCounter {
	return &RWCounter{
		mu:    sync.RWMutex{},
		count: 0,
	}
}

func (sc *SafeCounter) IncOne() {
	sc.mu.Lock()
	sc.count += 1
	sc.mu.Unlock()
}

func (sc *SafeCounter) DecOne() {
	sc.mu.Lock()
	sc.count -= 1
	sc.mu.Unlock()
}

func (sc *SafeCounter) Get() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}

func (rw *RWCounter) IncOne() {
	rw.mu.Lock()
	rw.count += 1
	rw.mu.Unlock()
}

func (rw *RWCounter) DecOne() {
	rw.mu.Lock()
	rw.count -= 1
	rw.mu.Unlock()
}

func (rw *RWCounter) Get() int {
	rw.mu.RLock()
	defer rw.mu.RUnlock()
	return rw.count
}
