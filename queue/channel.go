package queue

import "github.com/streadway/amqp"

type Channel struct {
	ServerUrl string
	LastError error
	*amqp.Connection
	*amqp.Channel
}

func NewChannel(url string) *Channel {
	if url == "" {
		panic("AMQP URL is empty !")
	}
	c := &Channel{ServerUrl: url}
	c.LastError = c.Reconnect(true)
	return c
}

func (c *Channel) Reconnect(force bool) (err error) {
	if force || c.Connection == nil {
		c.Close()
		c.Connection, err = amqp.Dial(c.ServerUrl)
	}
	c.Channel, err = c.Connection.Channel()
	return
}

func (c *Channel) Close() error {
	if c.Connection == nil {
		return nil
	}
	if c.Channel != nil {
		if err := c.Channel.Close(); err != nil {
			return err
		}
	}
	return c.Connection.Close()
}

func (c *Channel) Recover() error {
	if e := recover(); e != nil {
		c.LastError = e.(error)
	}
	return c.LastError
}

func (c *Channel) InitExchange(exchName, exchType string, durable bool) error {
	defer c.Recover()
	return c.ExchangeDeclare(
		exchName, // name
		exchType, // type
		durable,  // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // args
	)
}

func (c *Channel) InitQueue(queName string, durable bool) error {
	defer c.Recover()
	_, err := c.QueueDeclare(
		queName, // name
		durable, // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // args
	)
	return err
}

func (c *Channel) InitRoute(exchName, rtKey, queName string) error {
	defer c.Recover()
	return c.QueueBind(
		queName,  // queue
		rtKey,    // key
		exchName, // exchange
		false,    // no-wait
		nil,      // args
	)
}

func (c *Channel) InitBinds(exchName, exchType string, keys, ques []string) []error {
	// 忽略和已有定义不匹配的错误
	var (
		key     string
		err     error
		errs    []error
		durable = false // 性能较好
	)
	if err = c.InitExchange(exchName, exchType, durable); err != nil {
		errs = append(errs, err)
	}
	for i := 0; i < len(ques); i++ {
		if err = c.InitQueue(ques[i], durable); err != nil {
			errs = append(errs, err)
		}
		key = ""
		if i < len(keys) {
			key = keys[i]
		}
		if err = c.InitRoute(exchName, key, ques[i]); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (c *Channel) Redirect(src, rtKey, dst string) error {
	defer c.Recover()
	return c.ExchangeBind(
		dst,   // destination
		rtKey, // key
		src,   // source
		false, // no-wait
		nil,   // args
	)
}

func (c *Channel) RemoveQueue(queName string, isDel bool) (int, error) {
	defer c.Recover()
	if isDel {
		return c.QueueDelete(queName, false, false, false)
	} else {
		return c.QueuePurge(queName, false)
	}
}

func (c *Channel) ConsumeQueue(queName, csmTag string, autoAck bool) (<-chan amqp.Delivery, error) {
	defer c.Recover()
	return c.Consume(
		queName, // queue
		csmTag,  // consumer
		autoAck, // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
}

func (c *Channel) SetConfirm(confirm chan amqp.Confirmation) error {
	// publisher confirms for this channel/connection
	if err := c.Confirm(false); err != nil {
		close(confirm) // confirms not supported, simulate by always nacking
		return err
	} else {
		c.NotifyPublish(confirm)
		return nil
	}
}

func (c *Channel) PushMessage(exchName, key string, msg *Message) error {
	return c.Publish(
		exchName, // publish to an exchange
		key,      // routing to 0 or more queues
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Body:         msg.Body,
			Headers:      msg.Headers,
			DeliveryMode: 1, //非持久化
		},
	)
}
