package rdspool

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
)

type RedisHash struct {
	Inst Redis
	name string
	ttl  string
}

func NewRedisHash(inst Redis, name string, timeout int64) *RedisHash {
	ttl := strconv.FormatInt(timeout, 10)
	return &RedisHash{Inst: inst, name: name, ttl: ttl}
}

func (rh *RedisHash) DoWith(cmd string, args ...interface{}) (interface{}, error) {
	return DoWithKey(rh.Inst, cmd, rh.name, args...)
}

func (rh *RedisHash) Set(key, value string) (int, error) {
	reply, err := rh.Inst.Do("HSET", rh.name, key, value)
	if err == nil {
		rh.Inst.Do("SETTTL", rh.name, rh.ttl)
	}
	return redis.Int(reply, err)
}

func (rh *RedisHash) Get(key string) (string, error) {
	reply, err := rh.Inst.Do("HGET", rh.name, key)
	return redis.String(reply, err)
}

func (rh *RedisHash) GetAll() (map[string]string, error) {
	reply, err := rh.Inst.Do("HGETALL", rh.name)
	return redis.StringMap(reply, err)
}
