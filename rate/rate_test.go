package rate

import (
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

var (
	limiterSet       *Limiter
	limiterHset      *Limiter
	limiterRedisRate *redis_rate.Limiter
)

func init() {
	addr := "localhost:6379"
	pwd := "123456"
	limiterSet = NewLimiter(
		Addr(addr),
		Password(pwd),
		Window(WindowFixedSET),
	)

	limiterHset = NewLimiter(
		Addr(addr),
		Password(pwd),
		Window(WindowFixedHSET),
	)

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": addr,
		},
		Password: pwd,
	})
	limiterRedisRate = redis_rate.NewLimiter(ring)
	// Optional.
	limiterRedisRate.Fallback = rate.NewLimiter(rate.Every(time.Second), 100)
}

func TestLimiter_Set(t *testing.T) {
	for i := 0; i < 15; i++ {
		allocated, expiration, allow := limiterSet.AllowN("TestLimiter_Set", 10, time.Second, 1, "")
		t.Logf("allocated:%d, expiration:%v, allow:%v", allocated, expiration, allow)
	}
}

func TestLimiter_Hset(t *testing.T) {
	for i := 0; i < 15; i++ {
		allocated, expiration, allow := limiterHset.AllowN("TestLimiter_Hset", 10, time.Second, 1, "")
		t.Logf("allocated:%d, expiration:%v, allow:%v", allocated, expiration, allow)
	}
}

func TestLimiter_Rolling(t *testing.T) {

}

func TestLimiter_RedisRate(t *testing.T) {
	for i := 0; i < 15; i++ {
		c, d, a := limiterRedisRate.AllowN("TestLimiter_RedisRate", 10, time.Second, 1)
		t.Logf("count:%d, depay:%v, allow:%v", c, d, a)
	}
}

func BenchmarkLimiter_Set(b *testing.B) {
	b.Logf("N:%d", b.N)
	for i := 0; i < b.N; i++ {
		allocated, expiration, allow := limiterSet.AllowN("BenchmarkLimiter_Set", 5, time.Second, 1, "")
		b.Logf("allocated:%d, expiration:%v, allow:%v", allocated, expiration, allow)
	}
}

func BenchmarkLimiter_Hset(b *testing.B) {
	b.Logf("N:%d", b.N)
	for i := 0; i < b.N; i++ {
		allocated, expiration, allow := limiterHset.AllowN("BenchmarkLimiter_Hset", 5, time.Second, 1, "")
		b.Logf("allocated:%d, expiration:%v, allow:%v", allocated, expiration, allow)
	}
}

func BenchmarkLimiter_RedisRate(b *testing.B) {
	b.Logf("N:%d", b.N)
	for i := 0; i < b.N; i++ {
		c, d, a := limiterRedisRate.AllowN("BenchmarkLimiter_RedisRate", 5, time.Second, 1)
		b.Logf("count:%d, depay:%v, allow:%v", c, d, a)
	}
}
