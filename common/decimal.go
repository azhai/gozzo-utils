package common

import (
	"math"
	"strconv"
	"strings"
)

func RoundN(x float64, n int) float64 {
	multiple := math.Pow10(n)
	return math.Round(x*multiple) / multiple
}

// 高精度小数
// 例如 123.45 表示为 Decimal{Value:int64(12345), Precision:2}
type Decimal struct {
	Value     int64 // 扩大后成为整数
	Precision int   // 小数点后位数，限制15以内
}

// 使用方法：NewDecimal(123.45, 2)
func NewDecimal(value float64, prec int) *Decimal {
	d := &Decimal{}
	d.SetPrecision(prec)
	d.SetFloat(value, d.Precision)
	return d
}

// 使用方法：ParseDecimal("123.45", 2)
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
