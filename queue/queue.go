package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

type RecvFunc func(msg *Message) error
type TransFunc func(dlv amqp.Delivery) *Message

func FromDeliverie(dlv amqp.Delivery) *Message {
	if dlv.Body == nil {
		return nil
	}
	return &Message{Body: dlv.Body, Headers: dlv.Headers}
}

func IsValidateError(err error) bool {
	return strings.HasSuffix(err.Error(), " not supported")
}

// 消息队列
type MessageQueue struct {
	input    chan *Message
	output   <-chan amqp.Delivery
	QueName  string
	ExchName string
	Routings []string
}

func NewMessageQueue(queue, exchange string) *MessageQueue {
	return &MessageQueue{
		input:    make(chan *Message),
		QueName:  queue,
		ExchName: exchange,
	}
}

func (mq *MessageQueue) AddRoutings(key, dst string, count int) []string {
	var (
		routing, target string
		queNames        []string
	)
	for i := 0; i < count; i++ {
		if count <= 1 {
			routing = key
			target = dst
		} else {
			routing = fmt.Sprintf(key, i)
			target = fmt.Sprintf(dst, i)
		}
		queNames = append(queNames, target)
		mq.Routings = append(mq.Routings, routing)
	}
	return queNames
}

func (mq *MessageQueue) AddMessage(msg *Message) {
	mq.input <- msg
}

func (mq *MessageQueue) AddIndexMessage(msg *Message, index int) {
	if index >= 0 {
		msg.Routing = mq.Routings[index]
	}
	mq.AddMessage(msg)
}

func (mq *MessageQueue) AddData(body []byte, headers amqp.Table, index int) {
	msg := NewMessage(body)
	if headers != nil {
		for key, value := range headers {
			msg.Headers[key] = value
		}
	}
	mq.AddIndexMessage(msg, index)
}

func (mq *MessageQueue) PublishAll(ch *Channel, retries int) {
	go func() {
		var (
			err   error
			isVe  bool
			msg   *Message
			fails = 0
		)
		for {
			if err == nil || isVe {
				msg = <-mq.input
			}
			err = ch.PushMessage(mq.ExchName, msg.Routing, msg)
			if err == nil {
				continue
			}
			if isVe = IsValidateError(err); isVe {
				continue
			}
			if retries > 0 {
				if fails++; fails > retries {
					break
				}
			}
			time.Sleep(1 * time.Second)
			ch.Reconnect(true)
		}
	}()
}

func (mq *MessageQueue) NewTag(name string) string {
	return name + "-" + time.Now().Format("0102150405999")
}

func (mq *MessageQueue) Prepare(ch *Channel, queName string, retries int) (string, error) {
	var (
		err   error
		fails = 0
	)
	ctag := mq.NewTag(queName)
	mq.output, err = ch.ConsumeQueue(queName, ctag, true)
	for err != nil {
		time.Sleep(1 * time.Second)
		if ctag != "" {
			ch.Cancel(ctag, false)
			ctag = mq.NewTag(queName)
		}
		if err = ch.Reconnect(true); err == nil {
			mq.output, err = ch.ConsumeQueue(queName, ctag, true)
		}
		if retries > 0 {
			if fails++; fails > retries {
				break
			}
		}
	}
	return ctag, err
}

func (mq *MessageQueue) SubscribeAll(ch *Channel, receive RecvFunc) (err error) {
	if _, err = mq.Prepare(ch, mq.QueName, 3); err != nil {
		return
	}
	for dlv := range mq.output {
		if msg := FromDeliverie(dlv); msg != nil {
			go receive(msg)
		}
	}
	return
}

func (mq *MessageQueue) RunAll(ch *Channel, receive RecvFunc, retries int) {
	go func() {
		ctag, err := mq.Prepare(ch, mq.QueName, retries)
		if err != nil {
			return
		}
		for {
			select {
			case dlv := <-mq.output:
				if msg := FromDeliverie(dlv); msg != nil {
					go receive(msg)
				}
			case msg := <-mq.input:
				err := ch.PushMessage(mq.ExchName, msg.Routing, msg)
				if err != nil && !IsValidateError(err) {
					ctag, err = mq.Prepare(ch, ctag, retries)
				}
			}
		}
	}()
}
