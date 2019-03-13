package random

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const MAX_RAND_TIMES uint32 = 20000

var (
	// 使用随机种子初始化
	seedRand  = rand.New(rand.NewSource(time.Now().UnixNano()))
	randTimes = uint32(0)
)

// 2万次后更换source，因为C程序存在bug
// 参考 https://github.com/golang/go/issues/19809
// After 20k-50k iterations the loop starts printing i on every iteration.
func RandInt(n int) int {
	if atomic.AddUint32(&randTimes, 1) >= MAX_RAND_TIMES {
		atomic.StoreUint32(&randTimes, 0)
		seedRand.Seed(time.Now().UnixNano())
	}
	return seedRand.Intn(n)
}

// 从一定长度的数组中随机选取若干个索引，索引按正序（可能有循环）排列
func Sample(times, count int) (nums []int) {
	if count <= 0 || times <= 0 {
		return
	}
	if times < count {
		nums = JumpSample(times, count)
	} else {
		nums = PadSample(times, count)
	}
	return
}

// 数组长度足够，随机选取若干个索引，并保持正序
func JumpSample(times, count int) (nums []int) {
	// assert times < count
	step := count / times
	top := count%times + step
	i := -1
	for top <= count {
		i++
		if top <= i {
			break
		}
		i += RandInt(top - i)
		nums = append(nums, i)
		top += step
	}
	return
}

// 数组长度比要求的索引数量还少，按次序取出索引并不断翻倍
func PadSample(times, count int) (nums []int) {
	// assert times >= count
	for i := 1; i <= count; i++ {
		nums = append(nums, i)
	}
	var remain = times - count
	for remain > 0 {
		if remain > count {
			nums = append(nums, nums...)
			remain -= count
			count *= 2
		} else {
			nums = append(nums, nums[:remain]...)
			break
		}
	}
	return
}
