package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

type Envelope struct {
	QueueName   string
	Exchange    string // basic.publish exchange
	RoutingKey  string // basic.publish routing key
	*Message
}

type RecvFunc func(evp *Envelope) error

func FromDeliverie(queName string, dlv amqp.Delivery) *Envelope {
	if dlv.Body == nil {
		return nil
	}
	return &Envelope{
		QueueName: queName,
		Exchange: dlv.Exchange,
		RoutingKey: dlv.RoutingKey,
		Message: &Message{Body: dlv.Body, Headers: dlv.Headers},
	}
}

func IsValidateError(err error) bool {
	return strings.HasSuffix(err.Error(), " not supported")
}

// 消息队列
type MessageQueue struct {
	Input  chan *Envelope
	Handlers map[string]RecvFunc
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		Input:   make(chan *Envelope),
		Handlers: make(map[string]RecvFunc),
	}
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

func (mq *MessageQueue) AddMessage(exchName, routing string, msg *Message) {
	if msg != nil {
		evp := &Envelope{Exchange:exchName, RoutingKey:routing, Message: msg}
		mq.Input <- evp
	}
}

func (mq *MessageQueue) Build(exchName, routing string, body []byte, headers amqp.Table) {
	msg := NewMessage(body)
	if headers != nil {
		for key, value := range headers {
			msg.Headers[key] = value
		}
	}
	mq.AddMessage(exchName, routing, msg)
}

func (mq *MessageQueue) Publish(ch *Channel, input chan *Envelope, retries int) (err error) {
	fails, errch := 0, make(chan error, 1)
	var evp *Envelope
	for {
		select {
		case evp = <-input:
			err = ch.PushMessage(evp.Exchange, evp.RoutingKey, evp.Message)
			if err != nil && !IsValidateError(err) {
				errch <- err
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

func (mq *MessageQueue) Subscribe(ch *Channel, queName string, receive RecvFunc) (err error) {
	ctag := mq.NewTag(queName)
	defer ch.Cancel(ctag, false)
	var output <-chan amqp.Delivery
	expired := time.After(10 * time.Second) // 必须放在循环外
	for output == nil {
		select {
		case <-expired:
			return
		default:
			output, err = ch.ConsumeQueue(queName, ctag, true)
			if err != nil {
				ch.Close()
				time.Sleep(1 * time.Second)
				err = ch.Reconnect(false)
			}
		}
	}
	for dlv := range output {
		if evp := FromDeliverie(queName, dlv); evp != nil {
			go receive(evp)
		}
	}
	return
}

func (mq *MessageQueue) RunAll(ch *Channel, retries int) {
	for queName, receive := range mq.Handlers {
		go mq.Subscribe(ch, queName, receive)
	}
	pub := NewChannel(ch.ServerUrl)
	mq.Publish(pub, mq.Input, retries)
}
