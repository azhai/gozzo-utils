package redisw

import (
	"strconv"

	"github.com/azhai/gozzo-utils/common"
	"github.com/gomodule/redigo/redis"
)

const REDIS_DEFAULT_PORT = 6379

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
	return common.ConcatWith(defaultHost, p.StrPort(defaultPort))
}

func (p ConnParams) StrPort(defaultPort uint16) string {
	if p.Port > 0 {
		return strconv.Itoa(p.Port)
	}
	return strconv.Itoa(int(defaultPort))
}

func DialByParams(params ConnParams) (redis.Conn, error) {
	var opts []redis.DialOption
	addr := params.GetAddr("127.0.0.1", REDIS_DEFAULT_PORT)
	if params.Password != "" {
		opts = append(opts, redis.DialPassword(params.Password))
	}
	if dbno, err := strconv.Atoi(params.Database); err == nil {
		opts = append(opts, redis.DialDatabase(dbno))
	}
	return redis.Dial("tcp", addr, opts...)
}
