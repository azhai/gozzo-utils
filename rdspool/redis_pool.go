package rdspool

import (
	"github.com/azhai/gozzo-utils/common"
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
	pool := &RedisPool{addr: addr, passwd: passwd, db: db}
	pool.Reset(3, 5) //retryTimes=3, maxIdle=5
	return pool
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


func DoWith(r Redis, cmd, key string, args ...interface{}) (interface{}, error) {
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

// -1=无限 -2=不存在 -3=出错
func GetTimeout(r Redis, key string) int {
	sec, err := redis.Int(r.Do("TTL", key))
	if err == nil {
		return sec
	}
	return -3
}

func Expire(r Redis, key string, timeout int64) (bool, error) {
	reply, err := r.Do("EXPIRE", key, timeout)
	return redis.Bool(reply, err)
}

func Delete(r Redis, keys ...string) (int, error) {
	reply, err := r.Do("DEL", common.StrToList(keys)...)
	return redis.Int(reply, err)
}

func Exists(r Redis, key string) (bool, error) {
	reply, err := r.Do("EXISTS", key)
	return redis.Bool(reply, err)
}

func Find(r Redis, wildcard string) ([]string, error) {
	reply, err := r.Do("KEYS", wildcard)
	return redis.Strings(reply, err)
}

func Rename(r Redis, old, dst string) (bool, error) {
	reply, err := r.Do("RENAME", old, dst)
	return redis.Bool(reply, err)
}
