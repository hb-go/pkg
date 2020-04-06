package cache

import (
	"time"
)

var defaultCache Cache

func init() {
	// 默认go-cache
	defaultCache, _ = NewCacheGC()
}

func InitWithRedis(addrs string, password string) error {
	c, err := NewCacheRedis(addrs, password)
	if err != nil {
		return err
	} else {
		defaultCache = c
	}

	return nil
}

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string, value interface{}) error
	HSet(key string, field string, value interface{}) error
	HMSet(key string, fields map[string]interface{}) error
	HGet(key string, field string, value interface{}) error
	// HMGet() error
	// Add(key string, value interface{}, expire time.Duration) error
	// Replace(key string, data interface{}, expire time.Duration) error
	// IsExist(key string) bool
	// Delete(key string) error
	// Increment(key string, data uint64) (uint64, error)
	// Decrement(key string, data uint64) (uint64, error)
	// Flush() error
}

func Set(key string, value interface{}, ttl time.Duration) error {

	return defaultCache.Set(key, value, ttl)
}

func Get(key string, value interface{}) error {

	return defaultCache.Get(key, value)
}

func HSet(key string, field string, value interface{}) error {

	return defaultCache.HSet(key, field, value)
}

func HMSet(key string, fields map[string]interface{}) error {
	return defaultCache.HMSet(key, fields)
}

func HGet(key string, field string, value interface{}) error {

	return defaultCache.HGet(key, field, value)
}
