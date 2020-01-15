package redisw

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

type Map = map[string]interface{}

type ExecMulti func(keys ...string) (interface{}, error)

func ExecMap(exec ExecMulti, keys ...string) (data Map, err error) {
	var values []interface{}
	values, err = redis.Values(exec(keys...))
	if err != nil {
		return
	}
	data = make(Map, len(keys))
	for i, val := range values {
		data[keys[i]] = val
	}
	return
}

type CacheData interface {
	GetCacheId() string // 在同一个db中唯一id
}

////////////////////////////////////////////////////////////
/// redis string 的方法                                   ///
////////////////////////////////////////////////////////////

func (r *RedisWrapper) SetJson(key string, obj interface{}, timeout int64) (bool, error) {
	value, err := json.Marshal(obj)
	if err != nil {
		return false, err
	}
	return r.SetVal(key, value, timeout)
}

func (r *RedisWrapper) SetMap(data Map) (bool, error) {
	var args []interface{}
	for key, val := range data {
		args = append(args, key, val)
	}
	reply, err := r.Exec("MSET", args...)
	return redis.Bool(reply, err)
}

func (r *RedisWrapper) GetJson(key string, obj interface{}) error {
	value, err := r.GetBytes(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, &obj)
}

func (r *RedisWrapper) GetMulti(keys ...string) (interface{}, error) {
	args := StrToList(keys)
	return r.Exec("MGET", args...)
}

func (r *RedisWrapper) GetMap(keys ...string) (data Map, err error) {
	return ExecMap(r.GetMulti, keys...)
}

func (r *RedisWrapper) GetMapString(keys ...string) (map[string]string, error) {
	return redis.StringMap(r.GetMulti(keys...))
}

func (r *RedisWrapper) GetMapInt(keys ...string) (map[string]int, error) {
	return redis.IntMap(r.GetMulti(keys...))
}

////////////////////////////////////////////////////////////
/// redis hash 的方法                                     ///
////////////////////////////////////////////////////////////

func (rh *RedisHash) SetMap(data Map) (bool, error) {
	var args []interface{}
	for key, val := range data {
		args = append(args, key, val)
	}
	defer rh.Exec("EXPIRE", rh.timeout)
	return redis.Bool(rh.Exec("HMSET", args...))
}

func (rh *RedisHash) GetMulti(keys ...string) (interface{}, error) {
	args := StrToList(keys)
	return rh.Exec("HMGET", args...)
}

func (rh *RedisHash) GetMap(keys ...string) (data Map, err error) {
	return ExecMap(rh.GetMulti, keys...)
}

func (rh *RedisHash) GetMapString(keys ...string) (map[string]string, error) {
	return redis.StringMap(rh.GetMulti(keys...))
}

func (rh *RedisHash) GetMapInt(keys ...string) (map[string]int, error) {
	return redis.IntMap(rh.GetMulti(keys...))
}

////////////////////////////////////////////////////////////
/// redis string 和 hash 协作的方法                        ///
////////////////////////////////////////////////////////////

func (rh *RedisHash) SetMapJson(data Map) (bool, error) {
	summary, details := make(Map), make(Map)
	for key, val := range data {
		if obj, ok := val.(CacheData); ok {
			id := obj.GetCacheId()
			val, err := json.Marshal(val)
			if id != "" && err == nil {
				summary[key] = id
				details[id] = val
			}
		} else {
			summary[key] = val
		}
	}
	ok, err := rh.RedisWrapper.SetMap(details)
	//if ok && err == nil {
		ok, err = rh.SetMap(summary)
	//}
	return ok, err
}

func (rh *RedisHash) GetMapJson(data Map) error {
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	summary, err := rh.GetMapString(keys...)
	if err != nil {
		return err
	}
	for key, val := range summary {
		if _, ok := data[key].(CacheData); ok {
			err = rh.RedisWrapper.GetJson(val, data[key])
		} else {
			data[key] = val
		}
	}
	return err
}
