# gozzo 尜舟

## 常用 common

```go
package main
import (
    "fmt"
    "github.com/azhai/gozzo-utils/common"
)

// 浮点数
func main() {
    x := 123.45678
    a := common.NewDecimal(common.RoundN(x, 2), 2)
    fmt.Println(a.String()) // 123.45
    b := common.ParseDecimal(a.String(), 2)
    fmt.Println(b.String()) // 123.45
}
```

## 文件操作 filesystem

```go
package main
import (
    "fmt"
    "github.com/azhai/gozzo-utils/filesystem"
)

// 文件计行
func main() {
    fname := "README.md"
	count := LineCount(fname)

	// 逐行返回，适用于大文件
	var lines []string
	r := NewLineReader(fname)
	for r.Reading() {
		lines = append(lines, r.Text())
	}
	if len(lines) == count {
		fmt.Println("%s have %d lines", fname, count)
	} else {
		fmt.Println("Error !")
	}
}
```

## 文件日志 logging

```go
package main
import (
    "math"
    "time"
    "github.com/azhai/gozzo-utils/logging"
)

// 计算年龄
func CalcAge(birthday string) int {
    birth, err := time.Parse("2006-01-02", birthday)
    if err != nil {
        return -1
    }
    hours := time.Since(birth).Hours()
    return int(math.Round(hours / 365 / 24))
}

func main() {
    birthday := "1996-02-29"
    age := CalcAge(birthday)
    logger := logging.NewLogger("debug", "") // 输出到屏幕
    logger.Info("I was born on ", birthday, ", I am ",  age, " years old.")
}
```

## 地理位置和电子围栏 geohash

```bash
# 电子围栏的测试请查看文件 geohash/fence_test.go
cd geohash/
go test -v -mod=vendor
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
    routings := map[string]string{"testing":"queueForTesting"}
    ch.InitBinds("amq.topic", routings, true)
    mq := queue.NewMessageQueue()
    mq.AddHandler("queueForTesting", DumpBody) // 订阅队列testing
    mq.RunAll(ch, -1)
    for i := 1; i <= 10; i ++ {
        mq.AddMessage("amq.topic", "testing", CreateMessage(i)) // 发布消息
    }
}
```
