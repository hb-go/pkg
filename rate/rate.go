package rate

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	windowLUAScripts = map[WindowType]string{
		WindowFixedSET:  luaWindowFixedSet,
		WindowFixedHSET: luaWindowFixedHset,
		WindowRolling:   luaWindowRolling,
	}
)

type Limiter struct {
	client *redis.Client
	script *redis.Script

	opts *Options
}

func NewLimiter(opt ... Option) *Limiter {
	opts := DefaultOptions()
	for _, o := range opt {
		o(opts)
	}

	// initialize redis client
	option := redis.Options{
		Addr: opts.addr,
	}

	if len(opts.password) > 0 {
		option.Password = opts.password
	}

	if opts.poolSize > 0 {
		option.PoolSize = opts.poolSize
	}

	client := redis.NewClient(&option)
	script := redis.NewScript(windowLUAScripts[opts.window])

	return &Limiter{
		client: client,
		script: script,
		opts:   opts,
	}
}

func (l *Limiter) AllowN(key string, maxn int64, dur time.Duration, n int64, deduplicationID string) (allocated int64, expiration time.Duration, allow bool) {
	if l.opts.fallback != nil {
		allow = l.opts.fallback.Allow()
	}

	result, err := l.script.Run(
		l.client,
		[]string{
			key + ".meta", // KEY[1]
			key + ".data", // KEY[2]
		},
		// nolint: goimports
		maxn,                      // ARGV[1] credit
		dur.Nanoseconds(),         // ARGV[2] window length
		time.Second.Nanoseconds(), // ARGV[3] bucket length
		0,                         // ARGV[4] best effort
		n,                         // ARGV[5] token
		time.Now().UnixNano(),     // ARGV[6] timestamp
		deduplicationID,           // ARGS[7] deduplication id
	).Result()

	if err != nil {
		fmt.Printf("redis script run error:%v \n", err)
		return
	}

	allocated, expiration, err = getAllocatedTokenFromResult(&result)

	if err == nil {
		allow = allocated >= n
	}

	return
}

func getAllocatedTokenFromResult(result *interface{}) (int64, time.Duration, error) {
	if res, ok := (*result).([]interface{}); ok {
		if len(res) != 2 {
			return 0, 0, fmt.Errorf("invalid response from the redis server: %v", *result)
		}

		// read token
		tokenValue, tokenOk := res[0].(int64)
		if !tokenOk {
			return 0, 0, fmt.Errorf("invalid response from the redis server: %v", result)
		}

		// read expiration
		expValue, expOk := res[1].(int64)
		if !expOk {
			return 0, 0, fmt.Errorf("invalid response from the redis server: %v", result)
		}

		return tokenValue, time.Duration(expValue) * time.Nanosecond, nil
	}

	return 0, 0, fmt.Errorf("invalid response from the redis server: %v", result)
}
