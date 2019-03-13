# gozzo 尜舟

## 文件日志 logging

```go
package main
import (
    "github.com/azhai/gozzo-utils/common"
    "github.com/azhai/gozzo-utils/logging"
)

func main() {
    birthday := "1996-02-29"
    age := common.CalcAge(birthday)
    logger := logging.NewLogger("debug", "stderr", "stdout") // 输出到屏幕
    logger.Info("I was born on ", birthday, ", I am ",  age, " years old.")
}
```

## RabbitMQ队列 queue

```go
package main
import (
    "fmt"
    "github.com/streadway/amqp"
    "github.com/azhai/gozzo-utils/common"
    "github.com/azhai/gozzo-utils/queue"
)

// 创建JT/T808心跳消息，流水号为num
func CreateMessage(num int) *queue.Message {
    hb := common.Hex2Bin("7E0002000001453039919500")
    return &Message{
        Body: append(hb, 0x01 * num, 0x00, 0x7e), // 未计算校验码
        Headers: amqp.Table{
            "MsgNo": int16(num),
        },
        Routing: "testing", // 路由
    }
}

// 订阅消息的回调，直接输出消息体
func DumpBody(msg *queue.Message) error {
    fmt.Println(common.Bin2Hex(msg.Body))
    return nil
}

func main() {
    ch := queue.NewChannel("amqp://user:123@127.0.0.1:5672")
    defer ch.Close()
    mq := queue.NewMessageQueue("testing", "amq.topic") // 订阅队列testing
    mq.RunAll(ch, DumpBody, -1)
    for i := 1; i <= 10; i ++ {
        mq.AddMessage(CreateMessage(i)) // 发布消息
    }
}
```
