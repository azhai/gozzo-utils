package rdspool

import (
	"encoding/json"

	"github.com/azhai/gozzo-utils/common"
	"github.com/gomodule/redigo/redis"
)

type Map = map[string]interface{}

////////////////////////////////////////////////////////////
/// 以下为接口函数                                         ///
////////////////////////////////////////////////////////////

func SetNX(r Redis, key string, value interface{}) (int, error) {
	reply, err := r.Do("SETNX", key, value)
	return redis.Int(reply, err)
}

func SetVal(r Redis, key string, value interface{}, timeout int64) (bool, error) {
	reply, err := r.Do("SETEX", key, timeout, value)
	return redis.Bool(reply, err)
}

func SetPart(r Redis, key string, value interface{}, offset int) (int, error) {
	reply, err := r.Do("SETRANGE", key, offset, value)
	return redis.Int(reply, err)
}

func SetMap(r Redis, data Map) (bool, error) {
	var args []interface{}
	for key, val := range data {
		args = append(args, key, val)
	}
	reply, err := r.Do("MSET", args...)
	return redis.Bool(reply, err)
}

func SetJson(r Redis, key string, obj interface{}, timeout int64) (bool, error) {
	value, err := json.Marshal(obj)
	if err != nil {
		return false, err
	}
	return SetVal(r, key, value, timeout)
}

func GetBytes(r Redis, key string) ([]byte, error) {
	reply, err := r.Do("GET", key)
	return redis.Bytes(reply, err)
}

func GetString(r Redis, key string) (string, error) {
	reply, err := r.Do("GET", key)
	return redis.String(reply, err)
}

func GetStrLen(r Redis, key string) (int, error) {
	reply, err := r.Do("STRLEN", key)
	return redis.Int(reply, err)
}

func GetInt(r Redis, key string) (int, error) {
	reply, err := r.Do("GET", key)
	return redis.Int(reply, err)
}

// 计数增加
func IncrInt(r Redis, key string, offset int) (int, error) {
	value, err := IncrInt64(r, key, int64(offset))
	return int(value), err
}

func GetInt64(r Redis, key string) (int64, error) {
	reply, err := r.Do("GET", key)
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

func GetFloat(r Redis, key string) (float64, error) {
	reply, err := r.Do("GET", key)
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
	} else {
		reply, err = r.Do("INCRBYFLOAT", key, offset)
	}
	return redis.Float64(reply, err)
}

func GetMap(r Redis, keys ...string) (interface{}, error) {
	return r.Do("MGET", common.StrToList(keys)...)
}

func GetJson(r Redis, key string, obj interface{}) error {
	value, err := GetBytes(r, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, obj)
}
