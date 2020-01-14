package rdspool

import (
	"github.com/gomodule/redigo/redis"
)

type RedisHash struct {
	Inst    Redis
	name    string
	timeout int
}

func NewRedisHash(inst Redis, name string, timeout int) *RedisHash {
	return &RedisHash{Inst: inst, name: name, timeout: timeout}
}

func (rh *RedisHash) DoWith(cmd string, args ...interface{}) (interface{}, error) {
	return DoWithKey(rh.Inst, cmd, rh.name, args...)
}

// -1=无限 -2=不存在 -3=出错
func (rh *RedisHash) GetTimeout() int {
	if sec, err := redis.Int(rh.DoWith("TTL")); err == nil {
		return sec
	}
	return -3
}

func (rh *RedisHash) Set(key string, value interface{}) (int, error) {
	defer rh.DoWith("EXPIRE", rh.timeout)
	return redis.Int(rh.DoWith("HSET", key, value))
}

func (rh *RedisHash) Get(key string) (interface{}, error) {
	return rh.DoWith("HGET", key)
}

func (rh *RedisHash) GetString(key string) (string, error) {
	return redis.String(rh.Get(key))
}

func (rh *RedisHash) GetInt(key string) (int, error) {
	return redis.Int(rh.Get(key))
}

func (rh *RedisHash) GetInt64(key string) (int64, error) {
	return redis.Int64(rh.Get(key))
}

func (rh *RedisHash) GetFloat(key string) (float64, error) {
	return redis.Float64(rh.Get(key))
}

func (rh *RedisHash) GetAll() (interface{}, error) {
	return rh.DoWith("HGETALL")
}

func (rh *RedisHash) GetAllString(key string) (map[string]string, error) {
	return redis.StringMap(rh.GetAll())
}
