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
	for i, p := range pieces {
		pieces[i] = conv(p)
	}
	return pieces
}

// 高精度小数
type Decimal struct {
	Integer, Fraction int // 整数和小数部分
	Precision         int // 精确到小数点后第几位
}

func NewDecimal(value float64, prec int) Decimal {
	d := Decimal{Precision: prec, Integer: int(value)}
	remain := value - float64(d.Integer)
	remain *=  math.Pow10(d.Precision)
	d.Fraction = int(math.Round(remain))
	return d
}

func (d *Decimal) HasFraction() bool {
	return d.Precision > 0 && d.Fraction != 0
}

func (d *Decimal) CorrectPrecision(prec int) int {
	if prec > 9 {
		prec = 9
	} else if prec < 0 {
		prec = 0
	}
	return prec
}

func (d *Decimal) ChangePrecision(offset int) {
	oldPrec := d.Precision
	d.Precision = d.CorrectPrecision(d.Precision + offset)
	offset = d.Precision - oldPrec
	if offset == 0 || d.HasFraction() == false {
		return
	}
	if offset > 0 {
		d.Fraction *= int(math.Pow10(offset))
	} else {
		remain := float64(d.Fraction) / math.Pow10(offset)
		d.Fraction = int(math.Round(remain))
	}
}

func (d *Decimal) String() string {
	result := strconv.Itoa(d.Integer)
	if d.HasFraction() {
		tpl := "%0" + strconv.Itoa(d.Precision) + "d"
		frac := fmt.Sprintf(tpl, d.Fraction)
		result += "." + strings.TrimRight(frac, "0")
	}
	return result
}

// JSON中的日期时间类型
type JsonTime struct {
	Layout string // 格式，例如2006-01-02 15:04:05.999
	time.Time
}

func (t *JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.Format(t.Layout))
	return []byte(stamp), nil
}

func (t *JsonTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(t.Layout, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}
