package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 普通测试
func TestCommon(t *testing.T) {
	nums := Sample(13, 100)
	t.Log(nums)
	assert.Len(t, nums, 13)
}

// 测试短数组
func TestLess(t *testing.T) {
	nums := Sample(100, 13)
	t.Log(nums)
	assert.Len(t, nums, 100)
}

// 测试一样长时的数组
func TestEqual(t *testing.T) {
	nums := Sample(13, 13)
	t.Log(nums)
	assert.Len(t, nums, 13)
}

// 测试空数组
func TestEmpty(t *testing.T) {
	nums := Sample(13, 0)
	t.Log(nums)
	assert.Empty(t, nums)
}

// 测试非法参数
func TestEmpty2(t *testing.T) {
	nums := Sample(-1, 100)
	t.Log(nums)
	assert.Empty(t, nums)
}

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt(1<<32 - 1)
	}
}

func BenchmarkCommon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nums := JumpSample(13, 100)
		if size := len(nums); size != 13 {
			b.Errorf("The size of result %d is not 13", size)
		}
	}
}

func BenchmarkLess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nums := PadSample(100, 13)
		if size := len(nums); size != 100 {
			b.Errorf("The size of result %d is not 100", size)
		}
	}
}
