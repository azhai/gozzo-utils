package redisw

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

type CacheData interface {
	GetCacheId() string // 在同一个db中唯一id
}

type ExecMulti func(keys ...string) (interface{}, error)

type Map = map[string]interface{}

func NewMap() Map {
	return make(Map)
}

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

func ExecMapString(exec ExecMulti, keys ...string) (map[string]string, error) {
	values, err := redis.Strings(exec(keys...))
	if err != nil {
		return nil, err
	}
	data := make(map[string]string, len(values))
	for i, val := range values {
		data[keys[i]] = val
	}
	return data, err
}

func ExecMapBytes(exec ExecMulti, keys ...string) (map[string][]byte, error) {
	values, err := redis.ByteSlices(exec(keys...))
	if err != nil {
		return nil, err
	}
	data := make(map[string][]byte, len(values))
	for i, val := range values {
		data[keys[i]] = val
	}
	return data, err
}

func ExecMapInt(exec ExecMulti, keys ...string) (map[string]int, error) {
	values, err := redis.Ints(exec(keys...))
	if err != nil {
		return nil, err
	}
	data := make(map[string]int, len(values))
	for i, val := range values {
		data[keys[i]] = val
	}
	return data, err
}

func ExecMapInt64(exec ExecMulti, keys ...string) (map[string]int64, error) {
	values, err := redis.Int64s(exec(keys...))
	if err != nil {
		return nil, err
	}
	data := make(map[string]int64, len(values))
	for i, val := range values {
		data[keys[i]] = val
	}
	return data, err
}

func ExecMapFloat(exec ExecMulti, keys ...string) (map[string]float64, error) {
	values, err := redis.Float64s(exec(keys...))
	if err != nil {
		return nil, err
	}
	data := make(map[string]float64, len(values))
	for i, val := range values {
		data[keys[i]] = val
	}
	return data, err
}

////////////////////////////////////////////////////////////
/// redis string 的方法                                   ///
////////////////////////////////////////////////////////////

func (r *RedisWrapper) SaveJson(key string, obj interface{}, timeout int) (bool, error) {
	value, err := json.Marshal(obj)
	if err != nil {
		return false, err
	}
	return r.SetVal(key, value, timeout)
}

func (r *RedisWrapper) SaveMap(data Map) (bool, error) {
	var args []interface{}
	for key, val := range data {
		args = append(args, key, val)
	}
	reply, err := r.Exec("MSET", args...)
	return ReplyBool(reply, err)
}

func (r *RedisWrapper) LoadJson(key string, obj interface{}) error {
	value, err := r.GetBytes(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, obj)
}

func (r *RedisWrapper) GetMulti(keys ...string) (interface{}, error) {
	if len(keys) == 0 {
		return nil, KeysEmptyError
	}
	args := StrToList(keys)
	return r.Exec("MGET", args...)
}

func (r *RedisWrapper) LoadMap(keys ...string) (data Map, err error) {
	return ExecMap(r.GetMulti, keys...)
}

func (r *RedisWrapper) LoadMapString(keys ...string) (map[string]string, error) {
	return ExecMapString(r.GetMulti, keys...)
}

func (r *RedisWrapper) LoadMapInt(keys ...string) (map[string]int, error) {
	return ExecMapInt(r.GetMulti, keys...)
}

////////////////////////////////////////////////////////////
/// redis hash 的方法                                     ///
////////////////////////////////////////////////////////////

func (rh *RedisHash) SaveMap(data Map) (bool, error) {
	var args []interface{}
	for key, val := range data {
		args = append(args, key, val)
	}
	defer rh.Exec("EXPIRE", rh.timeout)
	return ReplyBool(rh.Exec("HMSET", args...))
}

func (rh *RedisHash) GetMulti(keys ...string) (interface{}, error) {
	if len(keys) == 0 {
		return nil, KeysEmptyError
	}
	args := StrToList(keys)
	return rh.Exec("HMGET", args...)
}

func (rh *RedisHash) LoadMap(keys ...string) (data Map, err error) {
	return ExecMap(rh.GetMulti, keys...)
}

func (rh *RedisHash) LoadMapString(keys ...string) (map[string]string, error) {
	return ExecMapString(rh.GetMulti, keys...)
}

func (rh *RedisHash) LoadMapInt(keys ...string) (map[string]int, error) {
	return ExecMapInt(rh.GetMulti, keys...)
}

////////////////////////////////////////////////////////////
/// redis string 和 hash 协作的方法                        ///
////////////////////////////////////////////////////////////

// 基本类型保存于自身，CacheData数据关联保存为Json
func (rh *RedisHash) SaveMapData(data Map) (ok bool, err error) {
	summary, timeout := NewMap(), rh.GetTimeout(true)
	for key, val := range data {
		if obj, ok := val.(CacheData); ok {
			id := obj.GetCacheId()
			if id != "" && err == nil {
				ok, err = rh.SaveJson(id, val, timeout)
				summary[key] = id
			}
		} else {
			summary[key] = val
		}
	}
	if err == nil {
		ok, err = rh.SaveMap(summary)
	}
	return
}

// 只能得到CacheData数据的Map，基本类型需要自己加载
func (rh *RedisHash) LoadMapJson(data Map) error {
	var keys []string
	for key, val := range data {
		if _, ok := val.(CacheData); !ok {
			continue
		}
		keys = append(keys, key)
	}
	summary, err := rh.LoadMapString(keys...)
	if err != nil {
		return err
	}
	for key, val := range summary {
		err = rh.RedisWrapper.LoadJson(val, data[key])
	}
	return err
}
