package queue

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

var (
	body     = []byte("*VK201867282020181224,AH&M510&N22&Z12&b2&T0000#")
	amqpUrl  = "amqp://user:123@192.168.2.107:5672"
	queName  = "TestQueue"
	prefix   = "mytestkey"
	messages = make([]*Message, 1000)
	counter  = int32(0)
)

func CreateMessage(tid int) *Message {
	tbin := []byte(strconv.Itoa(tid))
	copy(body[46-len(tbin):46], tbin[:])
	routing := prefix + "0"
	if tid%2 == 1 {
		routing = prefix + "1"
	}
	return &Message{
		Body: body,
		Headers: amqp.Table{
			"MsgId": int16(tid),
			"CmdId": "AH",
			"IMEI":  "867282020181224",
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
	keys := []string{prefix, prefix + "0", prefix + "1"}
	targets := []string{queName, queName + "0", queName + "1"}
	err := ch.InitBinds("amq.topic", "topic", keys, targets)
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
	ch.PushMessage("amq.topic", prefix, CreateMessage(sec))
	ch.PushMessage("amq.topic", prefix, CreateMessage(sec+1))
	ch.PushMessage("amq.topic", prefix, CreateMessage(sec+2))
}

func BenchmarkPublish1(b *testing.B) {
	GenMessages()
	ch := NewChannel(amqpUrl)
	defer ch.Close()
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	fmt.Println("amq.topic", prefix)
	for i := 0; i < b.N; i++ {
		idx := i % 1000
		ch.PushMessage("amq.topic", prefix, messages[idx])
	}
}

func BenchmarkPublish2(b *testing.B) {
	GenMessages()
	ch := NewChannel(amqpUrl)
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	mq := NewMessageQueue(queName, "amq.topic")
	targets := mq.AddRoutings(prefix+"%d", queName+"%d", 3)
	ch.InitBinds("amq.topic", "topic", mq.Routings, targets)
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
		t.Log(string(msg.Body))
		return nil
	}
}

func CreateCountFunc(b *testing.B) RecvFunc {
	return func(msg *Message) error {
		if len(msg.Body) < 47 {
			return fmt.Errorf("Too short body: %q", msg.Body)
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
	mq := NewMessageQueue(queName, "amq.topic")
	mq.RunAll(ch, CreateDumpFunc(t), -1)
}

func BenchmarkSubscribe(b *testing.B) {
	ch := NewChannel(amqpUrl)
	if ch.LastError != nil {
		b.Fatal(ch.LastError)
	}
	atomic.StoreInt32(&counter, 0)
	mq := NewMessageQueue(queName+"1", "amq.topic")
	mq.RunAll(ch, CreateCountFunc(b), -1)
	for i := 0; i < b.N; i += 2 {
		idx := i % 1000
		mq.AddMessage(messages[idx])
		mq.AddMessage(messages[idx+1])
	}
}
