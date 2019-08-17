package common

import (
	"fmt"
	"strings"
)

type WalkFunc func(item interface{}) error
type MapFunc func(item interface{}) (interface{}, error)
type ReduceFunc func(a, b interface{}) (interface{}, error)

// 将字符串数组转为一般数组
func StrToList(data []string) []interface{} {
	result := make([]interface{}, len(data))
	for i, v := range data {
		result[i] = v
	}
	return result
}

func SprintfString(tpl string, data []string) string {
	return fmt.Sprintf(tpl, StrToList(data)...)
}

func SprintfSplit(tpl string, data, sep string) string {
	return SprintfString(tpl, strings.Split(data, sep))
}

type IArray interface {
	ToList() []interface{}
}

// 循环修改数组
func ArrayWalk(arr IArray, f WalkFunc) error {
	for _, item := range arr.ToList() {
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

func ArrayMap(arr IArray, f MapFunc) ([]interface{}, error) {
	var (
		res []interface{}
		err error
	)
	for _, item := range arr.ToList() {
		if r, err := f(item); err != nil {
			res = append(res, r)
		}
	}
	return res, err
}

func ArrayReduce(arr IArray, f ReduceFunc, res interface{}) (interface{}, error) {
	var err error
	for _, item := range arr.ToList() {
		res, err = f(res, item)
	}
	return res, err
}

// 获取Slice的起止index
func GetStartStop(offset, limit, count int) (int, int) {
	if count < 0 || offset >= count || offset < 0-count || limit < 0-count {
		return -1, -1 // 参数不合理
	}
	if limit >= 0 && limit > count {
		limit = count
	}
	if offset < 0 {
		if limit < 0 {
			limit = limit - offset
			if limit <= 0 {
				return -1, -1
			}
		}
		offset = count + offset
	}
	start, stop := offset, offset+limit
	if limit <= 0 {
		stop = count + limit
	}
	if stop <= start {
		return start, start
	} else {
		return start, stop
	}
}
