package rdspool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetRedis() Redis {
	return NewRedisPool("127.0.0.1:6379", "", 0)
}

func TestInt(t *testing.T) {
	name := "test:a"
	r := GetRedis()
	SetVal(r, name, 39, 60)
	assert.Equal(t, 60, GetTimeout(r, name))
	a, err := GetInt(r, name)
	assert.NoError(t, err)
	assert.Equal(t, 39, a)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 58, GetTimeout(r, name))
}

func TestHash(t *testing.T) {
	name, key := "test:hash", "a"
	rh := NewRedisHash(GetRedis(), name, 2)
	rh.SetVal(key, 40)
	assert.Equal(t, 2, rh.GetTimeout())
	a, err := rh.GetInt(key)
	assert.NoError(t, err)
	assert.Equal(t, 40, a)
	time.Sleep(2 * time.Second)
	assert.Equal(t, -2, rh.GetTimeout())
}
