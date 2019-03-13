package common

import (
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

//16进制的字符表示和二进制字节表示互转
var Bin2Hex = hex.EncodeToString

var Hex2Bin = func(data string) []byte {
	if block, err := hex.DecodeString(data); err == nil {
		return block
	}
	return nil
}

//字符求余转为数字
func ToNum(c byte) uint8 {
	return uint8((c - '0' + 100) % 10)
}

//将数值或者数字转为字符串
func AsString(v interface{}) string {
	return fmt.Sprintf("%#v", v)
}

// 获取Slice的起止index
func GetStartStop(offset, limit, count int) (int, int) {
	if count < 0 || offset >= count || offset < 0-count {
		return -1, -1
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

// 计算年龄
func CalcAge(birthday string) int {
	birth, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return -1
	}
	hours := time.Since(birth).Hours()
	return int(math.Round(hours / 365 / 24))
}
