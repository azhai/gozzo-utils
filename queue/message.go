package queue

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/azhai/gozzo-utils/common"
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
	Body      []byte
	Headers   map[string]interface{}
	Exchange  string // basic.publish exchange
	Routing   string // basic.publish routing key
	QueueName string
}

func NewMessage(body []byte) *Message {
	return &Message{
		Body:    body,
		Headers: make(map[string]interface{}),
	}
}

func FromDeliverie(dlv amqp.Delivery) *Message {
	if dlv.Body == nil {
		return nil
	}
	return &Message{
		Body:     dlv.Body,
		Headers:  dlv.Headers,
		Exchange: dlv.Exchange,
		Routing:  dlv.RoutingKey,
	}
}

func (m *Message) GetHeaders() amqp.Table {
	return amqp.Table(m.Headers)
}

func (m *Message) GetHeaderDecimal(key string) amqp.Decimal {
	if value, ok := m.Headers[key]; ok {
		return value.(amqp.Decimal)
	}
	return amqp.Decimal{}
}

func (m *Message) GetHeaderInt(key string) int {
	if value, ok := m.Headers[key]; ok {
		return value.(int)
	}
	return 0
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

// 有可能是整数
func (m *Message) GetHeaderSafe(key string) string {
	if value, ok := m.Headers[key]; ok {
		switch v := value.(type) {
		case amqp.Decimal:
			d := common.Decimal{
				Value:     int64(v.Value),
				Precision: int(v.Scale),
			}
			return d.Format()
		case []byte:
			return common.Bin2Hex(v)
		case byte:
			return strconv.FormatInt(int64(v), 10)
		case int:
			return strconv.FormatInt(int64(v), 10)
		case int32:
			return strconv.FormatInt(int64(v), 10)
		case int16:
			return strconv.FormatInt(int64(v), 10)
		case int64:
			return strconv.FormatInt(v, 10)
		case string:
			return v
		default:
			return fmt.Sprintf("%+v", v)
		}
	}
	return ""
}

func (m *Message) GetHeaderString(key string) string {
	if value, ok := m.Headers[key]; ok {
		return value.(string)
	}
	return ""
}

func (m *Message) SetHeaderDecimal(key string, value int64, floats int) {
	val := int32(math.Round(float64(value) * math.Pow10(floats)))
	m.Headers[key] = amqp.Decimal{Value: val, Scale: uint8(floats)}
}

func (m *Message) SetHeaderInt64(key string, value int64, bits int) {
	if bits == 16 {
		m.Headers[key] = int16(value)
	} else if bits == 32 {
		m.Headers[key] = int32(value)
	} else if bits == 64 {
		m.Headers[key] = value
	} else {
		m.Headers[key] = int(value)
	}
}

func (m *Message) SetHeaderTime(key string, t time.Time) {
	m.Headers[key] = int64(t.UnixNano() / 1000000)
}

func (m *Message) SetHeaderTimeNow(key string) {
	m.SetHeaderTime(key, time.Now())
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

func (m *Message) ToString() string {
	return fmt.Sprintf("%s\t%s\t%+v\t%s", m.Exchange,
		m.Routing, m.Headers, common.Bin2Hex(m.Body))
}
