package queue

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/azhai/gozzo-utils/common"
	"github.com/streadway/amqp"
)

var (
	bodyTpl  = "7E00020000014530399195%04d0A7E"
	amqpUrl  = "amqp://user:123@127.0.0.1:5672"
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
		Body: common.Hex2Bin(body),
		Headers: amqp.Table{
			"MsgId": int16(tid),
			"CmdId": "0002",
			"IMEI":  "014530399195",
		},
		Routing: routing,
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
	keys := []string{prefix + "0", prefix + "1"}
	targets := []string{prefix + "0", prefix + "1"}
	err := ch.InitBinds("test.direct", "direct", keys, targets)
	if err != nil {
		t.Fatal(err)
		t.Log(keys)
	}
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
	ch.PushMessage("test.direct", msg.Routing, msg)
	msg = CreateMessage(sec + 1)
	ch.PushMessage("test.direct", msg.Routing, msg)
	msg = CreateMessage(sec + 2)
	ch.PushMessage("test.direct", msg.Routing, msg)
}

func BenchmarkPublish1(b *testing.B) {
	GenMessages()
	ch := NewChannel(amqpUrl)
	defer ch.Close()
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	fmt.Println("test.direct", prefix)
	for i := 0; i < b.N; i++ {
		idx := i % 1000
		msg := messages[idx]
		ch.PushMessage("test.direct", msg.Routing, msg)
	}
}

func BenchmarkPublish2(b *testing.B) {
	GenMessages()
	ch := NewChannel(amqpUrl)
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	mq := NewMessageQueue(queName, "test.direct")
	targets := mq.AddRoutings(prefix+"%d", prefix+"%d", 3)
	ch.InitBinds("test.direct", "direct", mq.Routings, targets)
	mq.PublishAll(ch, -1)
	fmt.Println(queName)
	for i := 0; i < b.N; i += 2 {
		idx := i % 1000
		mq.AddMessage(messages[idx])
		mq.AddMessage(messages[idx+1])
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
	mq := NewMessageQueue(queName, "test.direct")
	mq.RunAll(ch, CreateDumpFunc(t), -1)
}

func BenchmarkSubscribe(b *testing.B) {
	ch := NewChannel(amqpUrl)
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	atomic.StoreInt32(&counter, 0)
	mq := NewMessageQueue(queName+"1", "test.direct")
	mq.RunAll(ch, CreateCountFunc(b), -1)
	for i := 0; i < b.N; i += 2 {
		idx := i % 1000
		mq.AddMessage(messages[idx])
		mq.AddMessage(messages[idx+1])
	}
}
