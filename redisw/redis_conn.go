package redisw

import (
	"fmt"
	"strconv"
	"time"

	"github.com/azhai/gozzo-utils/common"
	"github.com/gomodule/redigo/redis"
	"github.com/gomodule/redigo/redisx"
)

var (
	StrToList = common.StrToList // 将字符串数组转为一般数组
	KeysEmptyError = fmt.Errorf("The param which named 'keys' must not empty !")
)

// redigo没有将应答中的OK转为bool值(2020-01-16)
func ReplyBool(reply interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	var answer string
	answer, err = redis.String(reply, err)
	return answer == "OK", err
}

// 用:号连接两个部分，如果后一部分也存在的话
func ConcatWith(master, slave string) string {
	if slave != "" {
		master += ":" + slave
	}
	return master
}

// Redis连接配置
type ConnParams struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Options  map[string]interface{}
}

func (p ConnParams) GetAddr(defaultHost string, defaultPort uint16) string {
	if p.Host != "" {
		defaultHost = p.Host
	}
	return ConcatWith(defaultHost, p.StrPort(defaultPort))
}

func (p ConnParams) StrPort(defaultPort uint16) string {
	if p.Port > 0 {
		return strconv.Itoa(p.Port)
	}
	return strconv.Itoa(int(defaultPort))
}

// Redis 容器，包括 *redis.Pool 和 *redisx.ConnMux 两个实现
type RedisContainer interface {
	Get() redis.Conn
	Close() error
}

// Redis
type RedisWrapper struct {
	MaxIdle     int
	MaxExecTime int // 命令最大执行时长
	RetryTimes  int // 重试次数
	RedisContainer
}

func DialByParams(params ConnParams) (redis.Conn, error) {
	var opts []redis.DialOption
	addr := params.GetAddr("127.0.0.1", 6379)
	if params.Password != "" {
		opts = append(opts, redis.DialPassword(params.Password))
	}
	if dbno, err := strconv.Atoi(params.Database); err == nil {
		opts = append(opts, redis.DialDatabase(dbno))
	}
	return redis.Dial("tcp", addr, opts...)
}

func NewRedisWrapper() *RedisWrapper {
	return &RedisWrapper{MaxIdle: 0, MaxExecTime: 5, RetryTimes: 3}
}

func NewRedisPool(params ConnParams, maxIdle int) *RedisWrapper {
	r := NewRedisWrapper()
	r.MaxIdle = maxIdle
	dial := func()(redis.Conn, error) { return DialByParams(params) }
	r.RedisContainer = &redis.Pool{Dial: dial, MaxIdle: r.MaxIdle}
	return r
}

func NewRedisConnMux(params ConnParams) *RedisWrapper {
	r := NewRedisWrapper()
	conn, _ := DialByParams(params)
	r.RedisContainer = redisx.NewConnMux(conn)
	return r
}

// 单命令最大执行时长
func (r *RedisWrapper) GetMaxExecDuration() time.Duration {
	if r.MaxExecTime > 0 {
		return time.Second * time.Duration(r.MaxExecTime)
	}
	return 0
}

// 执行命令，将会重试几次
func (r *RedisWrapper) Exec(cmd string, args ...interface{}) (interface{}, error) {
	var (
		err   error
		reply interface{}
	)
	med := r.GetMaxExecDuration()
	for i := 0; i < r.RetryTimes; i++ {
		if med > 0 {
			reply, err = redis.DoWithTimeout(r.Get(), med, cmd, args...)
		} else {
			reply, err = r.Get().Do(cmd, args...)
		}
		if err == nil {
			break
		}
	}
	return reply, err
}

func (r *RedisWrapper) GetSize() int {
	size, _ := redis.Int(r.Exec("DBSIZE"))
	return size
}

// -1=无限 -2=不存在 -3=出错
func (r *RedisWrapper) GetTimeout(key string) int {
	sec, err := redis.Int(r.Exec("TTL", key))
	if err == nil {
		return sec
	}
	return -3
}

func (r *RedisWrapper) Expire(key string, timeout int) (bool, error) {
	reply, err := r.Exec("EXPIRE", key, timeout)
	return ReplyBool(reply, err)
}

func (r *RedisWrapper) Delete(keys ...string) (int, error) {
	if len(keys) == 0 {
		return 0, KeysEmptyError
	}
	reply, err := r.Exec("DEL", StrToList(keys)...)
	return redis.Int(reply, err)
}

func (r *RedisWrapper) DeleteAll() (bool, error) {
	return ReplyBool(r.Exec("FLUSHDB"))
}

func (r *RedisWrapper) Exists(key string) (bool, error) {
	return ReplyBool(r.Exec("EXISTS", key))
}

func (r *RedisWrapper) Find(wildcard string) ([]string, error) {
	return redis.Strings(r.Exec("KEYS", wildcard))
}

func (r *RedisWrapper) Rename(old, dst string) (bool, error) {
	return ReplyBool(r.Exec("RENAME", old, dst))
}
