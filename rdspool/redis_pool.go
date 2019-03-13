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

func (this *RedisPool) Reset(retryTimes, maxIdle int) {
	this.retryTimes = retryTimes
	this.pool = redis.NewPool(this.Dial, maxIdle)
}

// 连接Redis
func (this *RedisPool) Dial() (redis.Conn, error) {
	opt := redis.DialDatabase(this.db)
	conn, err := redis.Dial("tcp", this.addr, opt)
	if err == nil && len(this.passwd) > 0 {
		conn.Do("AUTH", this.passwd)
	}
	return conn, err
}

// 从池中取出一个redis.Conn
func (this *RedisPool) Get() redis.Conn {
	if this.pool == nil {
		this.Reset(this.retryTimes, this.maxIdle)
	}
	return this.pool.Get()
}

// 关闭连接池和其中的连接
func (this *RedisPool) Close() error {
	if this.pool == nil {
		return nil
	}
	err := this.pool.Close()
	if err == nil {
		this.pool = nil
	}
	return err
}

// 执行命令，将会重试几次
func (this *RedisPool) Do(cmd string, args ...interface{}) (interface{}, error) {
	var (
		err   error
		reply interface{}
	)
	for i := 0; i < this.retryTimes; i++ {
		reply, err = this.Get().Do(cmd, args...)
		if err == nil {
			break
		}
	}
	return reply, err
}

////////////////////////////////////////////////////////////
/// 以下为接口函数                                         ///
////////////////////////////////////////////////////////////

func GetInt64(self Redis, key string) (int64, error) {
	reply, err := self.Do("GET", key)
	return redis.Int64(reply, err)
}

func SetInt64(self Redis, key string, value, timeout int64) (int64, error) {
	val := strconv.FormatInt(value, 10)
	ttl := strconv.FormatInt(timeout, 10)
	reply, err := self.Do("SET", key, val, ttl)
	return redis.Int64(reply, err)
}

// 计数增加
func IncrInt64(self Redis, key string, offset int64) (int64, error) {
	var (
		err   error
		reply interface{}
	)
	if offset == 0 {
		reply, err = self.Do("GET", key)
	} else if offset == 1 {
		reply, err = self.Do("INCR", key)
	} else {
		reply, err = self.Do("INCRBY", key, offset)
	}
	return redis.Int64(reply, err)
}
