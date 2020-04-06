package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/hb-go/pkg/log"
)

type redisCache struct {
	client *redis.Client
}

func (c *redisCache) Set(key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(key, value, ttl).Err()
}

func (c *redisCache) Get(key string, value interface{}) error {
	err := c.client.Get(key).Scan(value)
	if err == redis.Nil {
		log.Debugf("redis cache get nil, key: %s", key)
		return nil
	}

	return err
}

func (c *redisCache) HSet(key string, field string, value interface{}) error {
	return c.client.HSet(key, field, value).Err()
}

func (c *redisCache) HMSet(key string, fields map[string]interface{}) error {
	return c.client.HMSet(key, fields).Err()
}

func (c *redisCache) HGet(key string, field string, value interface{}) error {
	cmd := c.client.HGet(key, field)
	err := cmd.Scan(value)
	if err == redis.Nil {
		log.Infof("redis cache hget nil, key: %s, field: %s", key, field)
		return nil
	}

	return err
}

func NewCacheRedis(addrs string, password string) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	pong, err := client.Ping().Result()

	if err != nil || pong != "PONG" {
		panic(fmt.Sprintln("Redis ping,Is redis config error ? , error : ", err))
		return nil, err
	}
	return &redisCache{client: client}, nil
}
