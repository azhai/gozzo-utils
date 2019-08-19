package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ConvAction func(s string) string

// 找出其中的数字，不含负号和小数点
func GetNumber(data string) int64 {
	re := regexp.MustCompile("[0-9]+")
	data = re.FindString(data)
	num, err := strconv.ParseInt(data, 10, 64)
	if err == nil {
		return num
	}
	return -1
}

// 分拆为多个部分，并对每一段作处理
func SplitPieces(text, sep string, conv ConvAction) []string {
	pieces := strings.SplitN(text, sep, -1)
	if conv != nil {
		for i, p := range pieces {
			pieces[i] = conv(p)
		}
	}
	return pieces
}

// JSON中的日期时间类型
type JsonTime struct {
	time.Time
}

func (t JsonTime) GetLayout() string {
	return "2006-01-02 15:04:05"
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	l := t.GetLayout()
	stamp := fmt.Sprintf("\"%s\"", t.Format(l))
	return []byte(stamp), nil
}

func (t *JsonTime) UnmarshalJSON(buf []byte) error {
	l := t.GetLayout()
	tt, err := time.Parse(l, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

// 精确到毫秒,用于SqlServer等场景
type JsonTimeMS struct {
	time.Time
}

func (t JsonTimeMS) GetLayout() string {
	return "2006-01-02 15:04:05.999"
}
