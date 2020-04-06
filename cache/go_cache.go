package cache

import (
	"time"

	gc "github.com/patrickmn/go-cache"
)

type gcCache struct {
	c *gc.Cache
}

func (c *gcCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.c.Set(key, value, ttl)
	return nil
}

func (c *gcCache) Get(key string, value interface{}) error {
	if v, ok := c.c.Get(key); ok {
		value = v
	}

	return nil
}

func (c *gcCache) HSet(key string, field string, value interface{}) error {

	c.c.SetDefault(c.hKey(key, field), value)
	return nil
}

func (c *gcCache) HMSet(key string, fields map[string]interface{}) error {
	for k, v := range fields {
		c.HSet(key, k, v)
	}

	return nil
}

func (c *gcCache) HGet(key string, field string, value interface{}) error {
	if v, ok := c.c.Get(c.hKey(key, field)); ok {
		value = v
	}
	return nil
}

func (c *gcCache) hKey(key string, field string) string {
	return key + "_" + field

}

func NewCacheGC() (Cache, error) {
	c := gc.New(12*time.Hour, time.Hour)

	return &gcCache{c: c}, nil
}
