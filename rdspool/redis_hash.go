package rdspool

import (
	"github.com/azhai/gozzo-utils/common"
	"github.com/gomodule/redigo/redis"
)

type RedisHash struct {
	name    string
	timeout int
	Redis
}

func NewRedisHash(inst Redis, name string, timeout int) *RedisHash {
	return &RedisHash{Redis: inst, name: name, timeout: timeout}
}

func (rh *RedisHash) DoCmd(cmd string, args ...interface{}) (interface{}, error) {
	return DoWith(rh.Redis, cmd, rh.name, args...)
}

// -1=无限 -2=不存在 -3=出错
func (rh *RedisHash) GetTimeout() int {
	return GetTimeout(rh.Redis, rh.name)
}

func (rh *RedisHash) GetSize() int {
	size, _ := redis.Int(rh.DoCmd("HLEN"))
	return size
}

func (rh *RedisHash) GetKeys() []string {
	keys, _ := redis.Strings(rh.DoCmd("HKEYS"))
	return keys
}

func (rh *RedisHash) Delete(keys ...string) (int, error) {
	reply, err := rh.DoCmd("HDEL", common.StrToList(keys)...)
	return redis.Int(reply, err)
}

func (rh *RedisHash) Exists(key string) (bool, error) {
	return redis.Bool(rh.DoCmd("HEXISTS", key))
}

func (rh *RedisHash) SetNX(key string, value interface{}) (int, error) {
	affects, err := redis.Int(rh.DoCmd("HSETNX", key, value))
	if affects == 1 {
		rh.DoCmd("EXPIRE", rh.timeout)
	}
	return affects, err
}

func (rh *RedisHash) SetVal(key string, value interface{}) (int, error) {
	defer rh.DoCmd("EXPIRE", rh.timeout)
	return redis.Int(rh.DoCmd("HSET", key, value))
}

func (rh *RedisHash) SetMap(data Map) (bool, error) {
	var args []interface{}
	for key, val := range data {
		args = append(args, key, val)
	}
	defer rh.DoCmd("EXPIRE", rh.timeout)
	return redis.Bool(rh.DoCmd("HMSET", args...))
}

func (rh *RedisHash) GetVal(key string) (interface{}, error) {
	return rh.DoCmd("HGET", key)
}

func (rh *RedisHash) GetString(key string) (string, error) {
	return redis.String(rh.GetVal(key))
}

func (rh *RedisHash) GetInt(key string) (int, error) {
	return redis.Int(rh.GetVal(key))
}

func (rh *RedisHash) IncrInt(key string, offset int) (int, error) {
	value, err := rh.IncrInt64(key, int64(offset))
	return int(value), err
}

func (rh *RedisHash) GetInt64(key string) (int64, error) {
	return redis.Int64(rh.GetVal(key))
}

func (rh *RedisHash) IncrInt64(key string, offset int64) (int64, error) {
	return redis.Int64(rh.DoCmd("HINCRBY", key, offset))
}

func (rh *RedisHash) GetFloat(key string) (float64, error) {
	return redis.Float64(rh.GetVal(key))
}

func (rh *RedisHash) GetMap(keys ...string) (interface{}, error) {
	return rh.DoCmd("HMGET", common.StrToList(keys)...)
}

func (rh *RedisHash) GetMapString(keys ...string) (map[string]string, error) {
	return redis.StringMap(rh.GetMap(keys...))
}

func (rh *RedisHash) GetMapInt(keys ...string) (map[string]int, error) {
	return redis.IntMap(rh.GetMap(keys...))
}

func (rh *RedisHash) GetAll() (interface{}, error) {
	return rh.DoCmd("HGETALL")
}

func (rh *RedisHash) GetAllString() (map[string]string, error) {
	return redis.StringMap(rh.GetAll())
}
