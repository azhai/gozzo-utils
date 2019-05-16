package queue

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

// 消息头
func NewHeaders(protoId int) amqp.Table {
	return amqp.Table{
		"ProtocolId": int16(protoId),
		"IMEI":       "",
		"DeviceId":   "", // 非必需
		"SessId":     "",
		"CmdId":      "",       // 非必需
		"MsgId":      int16(0), // 非必需
		"RecvTime":   int64(0),
		"Type":       int16(0), // 普通/登录/退出/回应
		"Flags":      int16(0),
	}
}

// 消息
type Message struct {
	Body    []byte
	Headers amqp.Table
}

func NewMessage(body []byte) *Message {
	var headers = NewHeaders(0)
	return &Message{Headers: headers, Body: body}
}

func (m *Message) SetHeaderInt(key string, val int) {
	m.Headers[key] = int16(val)
}

func (m *Message) GetHeaderInt16(key string) int16 {
	if val, ok := m.Headers[key]; ok {
		return val.(int16)
	}
	return 0
}

func (m *Message) GetHeaderString(key string) string {
	if val, ok := m.Headers[key]; ok {
		return val.(string)
	}
	return ""
}

// 有可能是整数
func (m *Message) GetHeaderSafe(key string) string {
	if val, ok := m.Headers[key]; ok {
		switch v := val.(type) {
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

func (m *Message) GetCmdId() string {
	return m.GetHeaderString("CmdId")
}

func (m *Message) GetImei() string {
	return m.GetHeaderSafe("IMEI")
}

func (m *Message) ChangeProtoId(v int) {
	m.SetHeaderInt("ProtocolId", v)
}

func (m *Message) ChangeType(v int) {
	m.SetHeaderInt("Type", v)
}

func (m *Message) AddFlag(v int16, replace bool) {
	if !replace {
		if val := m.GetHeaderInt16("Flags"); val > 0 {
			v = v | val
		}
	}
	m.Headers["Flags"] = v
}

func (m *Message) SetRecvTime(t time.Time) {
	m.Headers["RecvTime"] = int64(t.UnixNano() / 1000000)
}

func (m *Message) SetRecvTimeNow() {
	m.SetRecvTime(time.Now())
}
