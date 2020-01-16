package redisw

import (
	"github.com/gomodule/redigo/redis"
)

func (r *RedisWrapper) SetNX(key string, value interface{}) (int, error) {
	reply, err := r.Exec("SETNX", key, value)
	return redis.Int(reply, err)
}

func (r *RedisWrapper) SetPart(key string, value interface{}, offset int) (int, error) {
	reply, err := r.Exec("SETRANGE", key, offset, value)
	return redis.Int(reply, err)
}

func (r *RedisWrapper) SetVal(key string, value interface{}, timeout int) (bool, error) {
	reply, err := r.Exec("SETEX", key, timeout, value)
	return ReplyBool(reply, err)
}

func (r *RedisWrapper) GetBytes(key string) ([]byte, error) {
	reply, err := r.Exec("GET", key)
	return redis.Bytes(reply, err)
}

func (r *RedisWrapper) GetString(key string) (string, error) {
	reply, err := r.Exec("GET", key)
	return redis.String(reply, err)
}

func (r *RedisWrapper) GetStrLen(key string) (int, error) {
	reply, err := r.Exec("STRLEN", key)
	return redis.Int(reply, err)
}

func (r *RedisWrapper) GetInt(key string) (int, error) {
	reply, err := r.Exec("GET", key)
	return redis.Int(reply, err)
}

// 计数增加
func (r *RedisWrapper) IncrInt(key string, offset int) (int, error) {
	value, err := r.IncrInt64(key, int64(offset))
	return int(value), err
}

func (r *RedisWrapper) GetInt64(key string) (int64, error) {
	reply, err := r.Exec("GET", key)
	return redis.Int64(reply, err)
}

// 计数增加
func (r *RedisWrapper) IncrInt64(key string, offset int64) (int64, error) {
	var (
		err   error
		reply interface{}
	)
	if offset == 0 {
		reply, err = r.Exec("GET", key)
	} else if offset == 1 {
		reply, err = r.Exec("INCR", key)
	} else {
		reply, err = r.Exec("INCRBY", key, offset)
	}
	return redis.Int64(reply, err)
}

func (r *RedisWrapper) GetFloat(key string) (float64, error) {
	reply, err := r.Exec("GET", key)
	return redis.Float64(reply, err)
}

// 计数增加
func (r *RedisWrapper) IncrFloat(key string, offset float64) (float64, error) {
	var (
		err   error
		reply interface{}
	)
	if offset == 0 {
		reply, err = r.Exec("GET", key)
	} else {
		reply, err = r.Exec("INCRBYFLOAT", key, offset)
	}
	return redis.Float64(reply, err)
}
