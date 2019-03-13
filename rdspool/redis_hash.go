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

func (this *RedisHash) Set(key, value string) (string, error) {
	reply, err := this.Inst.Do("HSET", this.name, key, value)
	if err == nil {
		this.Inst.Do("SETTTL", this.name, this.ttl)
	}
	return redis.String(reply, err)
}

func (this *RedisHash) Get(key string) (string, error) {
	reply, err := this.Inst.Do("HGET", this.name, key)
	return redis.String(reply, err)
}

func (this *RedisHash) GetAll() (map[string]string, error) {
	reply, err := this.Inst.Do("HGETALL", this.name)
	return redis.StringMap(reply, err)
}
