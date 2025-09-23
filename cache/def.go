package cache

import "time"

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Clear() error
	AcquireLock(key string, ttl time.Duration) (bool, error)
	ReleaseLock(key string) (bool, error)
}