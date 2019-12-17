package common

import (
	"fmt"
	"strings"
	"time"
)

const (
	LAYOUT_DATETIME       = "2006-01-02 15:04:05"
	LAYOUT_DATETIME_MILLS = "2006-01-02 15:04:05.999"
)

func ParseDate(layout, date string) (time.Time, error) {
	loc := time.Now().Location()
	return time.ParseInLocation(layout, date, loc)
}

func NewDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func ToDate(t time.Time) time.Time {
	year, month, day := t.Date()
	loc := t.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func Today() time.Time {
	return ToDate(time.Now())
}

// JSON中的日期时间类型
type JsonTime struct {
	time.Time
}

func (t JsonTime) GetLayout() string {
	return LAYOUT_DATETIME
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
	return LAYOUT_DATETIME_MILLS
}
