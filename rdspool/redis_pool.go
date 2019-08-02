package rdspool

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// Redis接口，包括redis.Conn和RedisPool两个实现
type Redis interface {
	Close() error
	Do(cmd string, args ...interface{}) (interface{}, error)
}

// Redis连接池
type RedisPool struct {
	pool       *redis.Pool
	addr       string
	passwd     string
	db         int
	retryTimes int
	maxIdle    int
}

func NewRedisPool(addr, passwd string, db int) *RedisPool {
	obj := &RedisPool{addr: addr, passwd: passwd, db: db}
	obj.Reset(3, 5) //retryTimes=3, maxIdle=5
	return obj
}

////////////////////////////////////////////////////////////
/// 以下为对象方法                                         ///
////////////////////////////////////////////////////////////

func (rp *RedisPool) Reset(retryTimes, maxIdle int) {
	rp.retryTimes = retryTimes
	rp.pool = &redis.Pool{Dial: rp.Dial, MaxIdle: maxIdle}
}

// 连接Redis
func (rp *RedisPool) Dial() (redis.Conn, error) {
	opt := redis.DialDatabase(rp.db)
	conn, err := redis.Dial("tcp", rp.addr, opt)
	if err == nil && len(rp.passwd) > 0 {
		conn.Do("AUTH", rp.passwd)
	}
	return conn, err
}

// 从池中取出一个redis.Conn
func (rp *RedisPool) Get() redis.Conn {
	if rp.pool == nil {
		rp.Reset(rp.retryTimes, rp.maxIdle)
	}
	return rp.pool.Get()
}

// 关闭连接池和其中的连接
func (rp *RedisPool) Close() error {
	if rp.pool == nil {
		return nil
	}
	err := rp.pool.Close()
	if err == nil {
		rp.pool = nil
	}
	return err
}

// 执行命令，将会重试几次
func (rp *RedisPool) Do(cmd string, args ...interface{}) (interface{}, error) {
	var (
		err   error
		reply interface{}
	)
	for i := 0; i < rp.retryTimes; i++ {
		reply, err = rp.Get().Do(cmd, args...)
		if err == nil {
			break
		}
	}
	return reply, err
}

////////////////////////////////////////////////////////////
/// 以下为接口函数                                         ///
////////////////////////////////////////////////////////////

func DoWithKey(r Redis, cmd, key string, args ...interface{}) (interface{}, error) {
	switch len(args) {
	case 0:
		return r.Do(cmd, key)
	case 1:
		return r.Do(cmd, key, args[0])
	case 2:
		return r.Do(cmd, key, args[0], args[1])
	default:
		args = append([]interface{}{key}, args...)
		return r.Do(cmd, args...)
	}
}

func GetInt64(r Redis, key string) (int64, error) {
	reply, err := r.Do("GET", key)
	return redis.Int64(reply, err)
}

func SetInt64(r Redis, key string, value, timeout int64) (int64, error) {
	val := strconv.FormatInt(value, 10)
	ttl := strconv.FormatInt(timeout, 10)
	reply, err := r.Do("SET", key, val, ttl)
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
