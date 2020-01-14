package rdspool

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
)

////////////////////////////////////////////////////////////
/// 以下为 int 接口函数                                    ///
////////////////////////////////////////////////////////////

func GetInt(r Redis, key string) (int, error) {
	reply, err := r.Do("GET", key)
	return redis.Int(reply, err)
}

func SetInt(r Redis, key string, value int, timeout int64) (int, error) {
	val := strconv.Itoa(value)
	ttl := strconv.FormatInt(timeout, 10)
	reply, err := r.Do("SETEX", key, ttl, val)
	return redis.Int(reply, err)
}

// 计数增加
func IncrInt(r Redis, key string, offset int) (int, error) {
	var (
		err   error
		reply interface{}
	)
	if offset == 0 {
		reply, err = r.Do("GET", key)
	} else if offset == 1 {
		reply, err = r.Do("INCR", key)
	} else {
		reply, err = r.Do("INCRBY", key, offset)
	}
	return redis.Int(reply, err)
}

////////////////////////////////////////////////////////////
/// 以下为 int64 接口函数                                  ///
////////////////////////////////////////////////////////////

func GetInt64(r Redis, key string) (int64, error) {
	reply, err := r.Do("GET", key)
	return redis.Int64(reply, err)
}

func SetInt64(r Redis, key string, value, timeout int64) (int64, error) {
	val := strconv.FormatInt(value, 10)
	ttl := strconv.FormatInt(timeout, 10)
	reply, err := r.Do("SETEX", key, ttl, val)
	return redis.Int64(reply, err)
}

// 计数增加
func IncrInt64(r Redis, key string, offset int64) (int64, error) {
	var (
		err   error
		reply interface{}
	)
	if offset == 0 {
		reply, err = r.Do("GET", key)
	} else if offset == 1 {
		reply, err = r.Do("INCR", key)
	} else {
		reply, err = r.Do("INCRBY", key, offset)
	}
	return redis.Int64(reply, err)
}

////////////////////////////////////////////////////////////
/// 以下为 float64 接口函数                                  ///
////////////////////////////////////////////////////////////

func GetFloat(r Redis, key string) (float64, error) {
	reply, err := r.Do("GET", key)
	return redis.Float64(reply, err)
}

func SetFloat(r Redis, key string, value float64, timeout int64) (float64, error) {
	val := strconv.FormatFloat(value, 'e', -1, 10)
	ttl := strconv.FormatInt(timeout, 10)
	reply, err := r.Do("SETEX", key, ttl, val)
	return redis.Float64(reply, err)
}

// 计数增加
func IncrFloat(r Redis, key string, offset float64) (float64, error) {
	var (
		err   error
		reply interface{}
	)
	if offset == 0 {
		reply, err = r.Do("GET", key)
	} else if offset == 1 {
		reply, err = r.Do("INCR", key)
	} else {
		reply, err = r.Do("INCRBY", key, offset)
	}
	return redis.Float64(reply, err)
}
