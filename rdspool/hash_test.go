package rdspool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetRedis() Redis {
	return NewRedisPool("127.0.0.1:6379", "", 0)
}

func GetRedisHash() *RedisHash {
	return NewRedisHash(GetRedis(), "test:hash", 2)
}

func TestStringHash(t *testing.T) {
	rh := GetRedisHash()
	rh.Set("a", 40)
	assert.Equal(t, 2, rh.GetTimeout())
	a, err := rh.GetInt("a")
	assert.NoError(t, err)
	assert.Equal(t, 40, a)
	time.Sleep(2 * time.Second)
	assert.Equal(t, -2, rh.GetTimeout())
}
