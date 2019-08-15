package common

import (
	"fmt"
	"math"
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

// 高精度小数
type Decimal struct {
	Value     int64 // 扩大后成为整数
	Precision int   // 小数点后位数，限制15以内
}

func NewDecimal(value float64, prec int) *Decimal {
	d := &Decimal{}
	d.SetPrecision(prec)
	d.SetFloat(value, d.Precision)
	return d
}

func ParseDecimal(text string, prec int) *Decimal {
	d := &Decimal{}
	d.SetPrecision(prec)
	if idx := strings.Index(text, "."); idx >= 0 {
		size := d.Precision + idx + 1
		if paddings := size - len(text); paddings > 0 {
			zeros := strings.Repeat("0", paddings)
			text = text[:idx] + text[idx+1:] + zeros
		} else {
			text = text[:idx] + text[idx+1:size]
		}
	}
	d.Value, _ = strconv.ParseInt(text, 10, 64)
	return d
}

func (d *Decimal) HasFraction() bool {
	if d.Precision <= 0 {
		return false
	}
	base := int64(math.Pow10(d.Precision))
	return d.Value%base != 0
}

func (d *Decimal) GetFloat() float64 {
	return float64(d.Value) / math.Pow10(d.Precision)
}

func (d *Decimal) SetFloat(value float64, expand int) {
	if expand > 0 {
		value *= math.Pow10(expand)
	}
	d.Value = int64(math.Round(value))
}

func (d *Decimal) SetPrecision(prec int) {
	if prec >= 15 {
		d.Precision = 15
	} else if prec <= 0 {
		d.Precision = 0
	} else {
		d.Precision = prec
	}
}

func (d *Decimal) ChangePrecision(offset int) {
	oldPrec := d.Precision
	d.SetPrecision(d.Precision + offset)
	offset = d.Precision - oldPrec
	if offset > 0 {
		d.Value *= int64(math.Pow10(offset))
	} else if offset < 0 {
		d.SetFloat(float64(d.Value), 0-offset)
	}
}

func (d *Decimal) String() string {
	result := strconv.FormatInt(d.Value, 10)
	if size := len(result) - d.Precision; size > 0 {
		result = result[:size] + "." + result[size:]
	} else {
		result = "0." + strings.Repeat("0", 0-size) + result
	}
	// 分开去除，否则会去掉整数部分末尾的0
	result = strings.TrimRight(result, "0")
	result = strings.TrimRight(result, ".")
	return result
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
