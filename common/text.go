package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func ReduceSpaces(lines string) string {
	return strings.Join(strings.Fields(lines), " ")
}

func WrapWith(s, left, right string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s%s%s", left, s, right)
}

func ReplaceWith(s string, subs map[string]string) string {
	if s == "" {
		return ""
	}
	var marks []string
	for key, value := range subs {
		marks = append(marks, key, value)
	}
	replacer := strings.NewReplacer(marks...)
	return replacer.Replace(s)
}

func ReplaceQuotes(s string) string {
	if s == "" {
		return ""
	}
	replacer := strings.NewReplacer("[", "`", "]", "`")
	return replacer.Replace(s)
}