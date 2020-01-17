package redisw

import (
	"encoding/json"
	"strings"

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

func Map2Args(data Map, asJson bool) []interface{} {
	var args []interface{}
	for key, val := range data {
		if asJson {
			value, err := json.Marshal(val)
			if err == nil {
				args = append(args, key, value)
			}
		} else if val != nil {
			args = append(args, key, val)
		}
	}
	return args
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

func (r *RedisWrapper) SaveMap(data Map, asJson bool) (bool, error) {
	args := Map2Args(data, asJson)
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

func (rh *RedisHash) SaveMap(data Map, asJson bool) (bool, error) {
	args := Map2Args(data, asJson)
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
func (rh *RedisHash) SaveForeignData(data Map) (ok bool, err error) {
	summary, timeout := NewMap(), rh.GetTimeout(true)
	for key, val := range data {
		if val == nil {
			continue
		}
		if obj, ok := val.(CacheData); ok {
			id := obj.GetCacheId()
			if id != "" && err == nil {
				ok, err = rh.RedisWrapper.SaveJson(id, val, timeout)
				summary[key] = id
			}
		} else {
			summary[key] = val
		}
	}
	if err == nil {
		ok, err = rh.SaveMap(summary, false)
	}
	return
}

func (rh *RedisHash) LoadSummary(data Map) (map[string]string, error) {
	var keys []string
	for key, val := range data {
		if val == nil {
			keys = append(keys, key)
		} else if _, ok := val.(CacheData); ok {
			keys = append(keys, key)
		}
	}
	return rh.LoadMapString(keys...)
}

// 只能得到CacheData数据的Map，基本类型需要自己加载
func (rh *RedisHash) LoadForeignJson(data Map) (err error) {
	var summary map[string]string
	if summary, err = rh.LoadSummary(data); err != nil {
		return err
	}
	for key, val := range summary {
		err = rh.RedisWrapper.LoadJson(val, data[key])
	}
	return err
}

func (rh *RedisHash) LoadForeignString(keys ...string) (result map[string]string, err error) {
	var summary, foreigns map[string]string
	if summary, err = rh.LoadMapString(keys...); err != nil {
		return
	}
	var ids []string
	for _, id := range summary {
		ids = append(ids, id)
	}
	if foreigns, err = rh.RedisWrapper.LoadMapString(ids...); err != nil {
		return
	}
	result = make(map[string]string)
	for key, id := range summary {
		if val, ok := foreigns[id]; ok {
			result[key] = val
		}
	}
	return
}

// 对list、dict或string的json，去掉两边括号或引号取中间部分
// 用于快速拼接json，免除判断类型和空值，如 fmt.Sprintf("{%s}", GetJsonContent(json))
func GetJsonContent(data string) string {
	data = strings.TrimSpace(data)
	if len(data) < 2 {
		return ""
	}
	if strings.HasPrefix(data, "{") && strings.HasSuffix(data, "}") {
		return data[1 : len(data)-1]
	}
	if strings.HasPrefix(data, "[") && strings.HasSuffix(data, "]") {
		return data[1 : len(data)-1]
	}
	if strings.HasPrefix(data, "\"") && strings.HasSuffix(data, "\"") {
		return data[1 : len(data)-1]
	}
	return ""
}