package queue

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/azhai/gozzo-utils/common"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

var (
	bodyTpl  = "7E00020000014530399195%04d0A7E"
	amqpUrl  = "amqp://guest:guest@127.0.0.1:5672"
	exchName = "amq.topic"
	queName  = "TestQueue"
	prefix   = "TestRouting"
	messages = make([]*Message, 1000)
	counter  = int32(0)
)

func CreateMessage(tid int) *Message {
	body := fmt.Sprintf(bodyTpl, tid)
	routing := prefix + "0"
	if tid%2 == 1 {
		routing = prefix + "1"
	}
	return &Message{
		Exchange: exchName,
		Routing:  routing,
		Body:     common.Hex2Bin(body),
		Headers: amqp.Table{
			"MsgId": int16(tid),
			"CmdId": "0002",
			"IMEI":  "014530399195",
		},
	}
}

func GenMessages() {
	for i := 0; i < cap(messages); i++ {
		messages[i] = CreateMessage(i + 1)
	}
}

func TestCreate(t *testing.T) {
	ch := NewChannel(amqpUrl)
	defer ch.Close()
	if ch.LastError != nil {
		t.Fatal(ch.LastError)
		return
	}
	if err := ch.InitQueue(queName, false); err != nil {
		t.Fatal(err)
	}
	routingMap := make(map[string]string)
	routingMap[prefix+"0"] = prefix + "0"
	routingMap[prefix+"1"] = prefix + "1"
	err := ch.InitBinds(exchName, routingMap, true)
	assert.NoError(t, err)
	t.Log(routingMap)
}

// 发布消息
func TestPublish(t *testing.T) {
	ch := NewChannel(amqpUrl)
	defer ch.Close()
	if ch.LastError != nil {
		t.Fatal(ch.LastError)
	}
	sec := time.Now().Second()
	msg := CreateMessage(sec)
	err := ch.InitExchange(msg.Exchange, "topic", true)
	assert.NoError(t, err)
	err = ch.PushMessage(msg)
	assert.NoError(t, err)
	msg = CreateMessage(sec + 1)
	err = ch.PushMessage(msg)
	assert.NoError(t, err)
	msg = CreateMessage(sec + 2)
	err = ch.PushMessage(msg)
	assert.NoError(t, err)
}

func BenchmarkPublish1(b *testing.B) {
	GenMessages()
	ch := NewChannel(amqpUrl)
	defer ch.Close()
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	fmt.Println("test.direct", prefix)
	_ = ch.InitExchange("test.direct", "direct", false)
	for i := 0; i < b.N; i++ {
		idx := i % 1000
		_ = ch.PushMessage(messages[idx])
	}
}

func BenchmarkPublish2(b *testing.B) {
	GenMessages()
	ch := NewChannel(amqpUrl)
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	mq := NewMessageQueue()
	routingMap := mq.AddRoutings(prefix+"%d", prefix+"%d", 3)
	_ = ch.InitBinds(exchName, routingMap, true)
	_ = mq.Publish(ch, mq.Input, -1)
	for i := 0; i < b.N; i += 2 {
		idx := i % 1000
		mq.Input <- messages[idx]
		mq.Input <- messages[idx+1]
	}
}

func CreateDumpFunc(t *testing.T) RecvFunc {
	return func(msg *Message) error {
		t.Log(common.Bin2Hex(msg.Body))
		return nil
	}
}

func CreateCountFunc(b *testing.B) RecvFunc {
	return func(msg *Message) error {
		if len(msg.Body) < 15 { // 最少15个字节
			return fmt.Errorf("Too short body: %d bytes", len(msg.Body))
		}
		val := atomic.AddInt32(&counter, 1)
		if val%20000 == 1 {
			b.Log(val)
		}
		return nil
	}
}

// 订阅消息
func TestSubscribe(t *testing.T) {
	ch := NewChannel(amqpUrl)
	if ch.LastError != nil {
		t.Fatal(ch.LastError)
	}
	mq := NewMessageQueue()
	mq.AddHandler(queName, CreateDumpFunc(t))
	err := mq.RunAll(ch, -1)
	assert.NoError(t, err)
}

func BenchmarkSubscribe(b *testing.B) {
	ch := NewChannel(amqpUrl)
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	atomic.StoreInt32(&counter, 0)
	mq := NewMessageQueue()
	mq.AddHandler(queName+"1", CreateCountFunc(b))
	_ = mq.RunAll(ch, -1)
	for i := 0; i < b.N; i += 2 {
		idx := i % 1000
		mq.Input <- messages[idx]
		mq.Input <- messages[idx+1]
	}
}
