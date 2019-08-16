package queue

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

/*
amqp.Table stores user supplied fields of the following types:
bool
byte
float32
float64
int
int16
int32
int64
nil
string
time.Time
amqp.Decimal
amqp.Table
[]byte
[]interface{} - containing above types
*/

// 消息
type Message struct {
	Body    []byte
	Headers amqp.Table
}

func NewMessage(body []byte) *Message {
	return &Message{Body: body, Headers: amqp.Table{}}
}

func (m *Message) GetHeaderString(key string) string {
	if value, ok := m.Headers[key]; ok {
		return value.(string)
	}
	return ""
}

// 有可能是整数
func (m *Message) GetHeaderSafe(key string) string {
	if value, ok := m.Headers[key]; ok {
		switch v := value.(type) {
		case string:
			return v
		case []byte:
			return hex.EncodeToString(v)
		case int64:
			return strconv.FormatInt(v, 10)
		case int32:
			return strconv.FormatInt(int64(v), 10)
		case int16:
			return strconv.FormatInt(int64(v), 10)
		default:
			return fmt.Sprintf("%+v", v)
		}
	}
	return ""
}

func (m *Message) SetHeaderInt(key string, value, bits int) {
	if bits == 16 {
		m.Headers[key] = int16(value)
	} else if bits == 32 {
		m.Headers[key] = int32(value)
	} else {
		m.Headers[key] = int64(value)
	}
}

func (m *Message) GetHeaderInt16(key string) int16 {
	if value, ok := m.Headers[key]; ok {
		return value.(int16)
	}
	return 0
}

func (m *Message) GetHeaderInt32(key string) int32 {
	if value, ok := m.Headers[key]; ok {
		return value.(int32)
	}
	return 0
}

func (m *Message) GetHeaderInt64(key string) int64 {
	if value, ok := m.Headers[key]; ok {
		return value.(int64)
	}
	return 0
}

func (m *Message) AddFlag(key string, v int16, replace bool) int16 {
	if !replace {
		if value := m.GetHeaderInt16(key); value > 0 {
			v = v | value
		}
	}
	m.Headers[key] = v
	return v
}

func (m *Message) SetTime(key string, t time.Time) {
	m.Headers[key] = int64(t.UnixNano() / 1000000)
}

func (m *Message) SetTimeNow(key string) {
	m.SetTime(key, time.Now())
}
