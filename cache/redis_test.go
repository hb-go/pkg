package cache

import (
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	c, _ := NewCacheRedis("localhost:6379", "123456")

	if err := c.Set("key", "value", time.Second*5); err != nil {
		t.Error(err)
	}
	v := ""
	if err := c.Get("key", &v); err != nil {
		t.Error(err)
	}
	t.Logf("value: %v", v)

	if err := c.HSet("hkey", "field", "value"); err != nil {
		t.Error(err)
	}
	hv := ""
	if err := c.HGet("hkey", "field", &hv); err != nil {
		t.Error(err)
	}
	t.Logf("h value: %v", hv)
}
