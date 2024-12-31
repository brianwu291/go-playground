package ratelimiter

import (
	"sync"
	"time"
)

type bucket struct {
	count    int
	mu       sync.RWMutex
	lastTime time.Time
}

type RateLimiter struct {
	duration time.Duration
	max      int
	buckets  map[string]*bucket
	mu       sync.RWMutex
}

const defaultDuration = time.Minute * 3

func NewRateLimiter(max int, duration time.Duration) *RateLimiter {
	if duration == 0 {
		duration = defaultDuration
	}
	return &RateLimiter{
		duration: duration,
		max:      max,
		buckets:  make(map[string]*bucket),
	}
}

func (r *RateLimiter) getBucketByKey(key string) *bucket {
	r.mu.Lock()
	defer r.mu.Unlock()

	b, exists := r.buckets[key]
	if !exists {
		b = &bucket{
			lastTime: time.Now(),
		}
		r.buckets[key] = b
	}
	return b
}

func (r *RateLimiter) Access(key string) bool {
	b := r.getBucketByKey(key)
	b.mu.Lock()
	defer b.mu.Unlock()

	nextResetTime := b.lastTime.Add(r.duration)
	if time.Until(nextResetTime) <= 0 {
		b.count = 0
		b.lastTime = time.Now()
	}

	if b.count >= r.max {
		return false
	}

	b.count += 1
	return true
}

func (r *RateLimiter) GetRestTime(key string) time.Duration {
	b := r.getBucketByKey(key)
	b.mu.RLock()
	defer b.mu.RUnlock()

	nextResetTime := b.lastTime.Add(r.duration)
	remaining := time.Until(nextResetTime)
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (r *RateLimiter) Clean() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, b := range r.buckets {
		b.mu.Lock()
		if time.Until(b.lastTime.Add(r.duration)) <= 0 {
			delete(r.buckets, key)
		}
		b.mu.Unlock()
	}
}
