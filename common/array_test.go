package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetElements(offset, limit int) []int {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	count := len(data)
	start, stop := GetStartStop(offset, limit, count)
	if start >= 0 && stop >= 0 {
		return data[start:stop]
	}
	return nil
}

func TestListSlice(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	t.Log(0, 0, GetElements(0, 0))
	assert.Equal(t, data, GetElements(0, 0))
	t.Log(0, 1, GetElements(0, 1))
	assert.Equal(t, []int{1}, GetElements(0, 1))
	t.Log(0, -1, GetElements(0, -1))
	assert.Equal(t, data[:8], GetElements(0, -1))
	t.Log(1, 0, GetElements(1, 0))
	assert.Equal(t, data[1:], GetElements(1, 0))
	t.Log(1, 1, GetElements(1, 1))
	assert.Equal(t, []int{2}, GetElements(1, 1))
	t.Log(1, -1, GetElements(1, -1))
	assert.Equal(t, data[1:8], GetElements(1, -1))
	t.Log(-1, 0, GetElements(-1, 0))
	assert.Equal(t, []int{9}, GetElements(-1, 0))
	t.Log(-1, 1, GetElements(-1, 1))
	assert.Equal(t, []int{9}, GetElements(-1, 1))
	t.Log(-1, -1, GetElements(-1, -1))
	assert.Empty(t, GetElements(-1, -1))
}
