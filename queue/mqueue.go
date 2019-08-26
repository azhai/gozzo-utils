package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/azhai/gozzo-utils/logging"
	"github.com/streadway/amqp"
)

type RecvFunc func(msg *Message) error

func IsValidateError(err error) bool {
	return strings.HasSuffix(err.Error(), " not supported")
}

// 消息队列
type MessageQueue struct {
	logger   logging.ILogger
	Input    chan *Message
	Handlers map[string]RecvFunc
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		Input:    make(chan *Message),
		Handlers: make(map[string]RecvFunc),
	}
}

func (mq *MessageQueue) SetLogger(logger logging.ILogger) {
	mq.logger = logger
}

func (mq *MessageQueue) AddRoutings(key, dst string, count int) map[string]string {
	var routing, target string
	routingMap := make(map[string]string)
	for i := 0; i < count; i++ {
		if count <= 1 {
			routing = key
			target = dst
		} else {
			routing = fmt.Sprintf(key, i)
			target = fmt.Sprintf(dst, i)
		}
		routingMap[target] = routing
	}
	return routingMap
}

func (mq *MessageQueue) AddMessage(msg *Message, routExch ...string) {
	if msg == nil {
		return
	}
	if size := len(routExch); size >= 1 {
		msg.Routing = routExch[0]
		if size == 2 {
			msg.Exchange = routExch[1]
		}
	}
	mq.Input <- msg
	if mq.logger != nil {
		mq.logger.Debug(msg.ToString())
	}
}

func (mq *MessageQueue) Build(exch, routing string, body []byte, headers amqp.Table) {
	msg := NewMessage(body)
	if headers != nil {
		for key, value := range headers {
			msg.Headers[key] = value
		}
	}
	mq.AddMessage(msg, exch, routing)
}

func (mq *MessageQueue) Publish(ch *Channel, input chan *Message, retries int) (err error) {
	fails, errch := 0, make(chan error, 1)
	var msg *Message
	for {
		select {
		case msg = <-input:
			err = ch.PushMessage(msg)
			if err != nil {
				if !IsValidateError(err) {
					errch <- err
				}
				if mq.logger != nil {
					mq.logger.Error(err.Error())
				}
			}
		case err = <-errch:
			if retries > 0 {
				if fails++; fails > retries {
					break
				}
			}
			ch.Close()
			time.Sleep(1 * time.Second)
			if err = ch.Reconnect(false); err != nil {
				errch <- err
				if mq.logger != nil {
					mq.logger.Error(err.Error())
				}
			}
		}
	}
	err = <-errch
	return
}

func (mq *MessageQueue) AddHandler(queName string, receive RecvFunc) {
	mq.Handlers[queName] = receive
}

func (mq *MessageQueue) NewTag(name string) string {
	return name + "-" + time.Now().Format("0102150405999")
}

func (mq *MessageQueue) Subscribe(ch *Channel, queName string, autoAck bool, receive RecvFunc) (err error) {
	ctag := mq.NewTag(queName)
	defer ch.Cancel(ctag, false)
	var output <-chan amqp.Delivery
	expired := time.After(10 * time.Second) // 必须放在循环外
	for output == nil {
		select {
		case <-expired:
			return
		default:
			output, err = ch.ConsumeQueue(queName, ctag, autoAck)
			if err != nil {
				ch.Close()
				if mq.logger != nil {
					mq.logger.Error(err.Error())
				}
				time.Sleep(1 * time.Second)
				err = ch.Reconnect(false)
			}
		}
	}
	for dlv := range output {
		if msg := FromDeliverie(dlv); msg != nil {
			msg.QueueName = queName
			go receive(msg)
		}
	}
	return
}

func (mq *MessageQueue) RunAll(ch *Channel, retries int) (err error) {
	if ch.LastError != nil {
		if err = ch.Reconnect(true); err != nil {
			if mq.logger != nil {
				mq.logger.Error(err.Error())
			}
			ch.LastError = err
			return
		}
	}
	for queName, receive := range mq.Handlers {
		go mq.Subscribe(ch, queName, true, receive)
	}
	pub := NewChannel(ch.ServerUrl)
	return mq.Publish(pub, mq.Input, retries)
}
