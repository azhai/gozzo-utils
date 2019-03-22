package common

import "testing"

func GetElements(offset, limit int) []int {
	var data []int
	for i := 1; i <= 9; i++ {
		data = append(data, i)
	}
	count := len(data)
	start, stop := GetStartStop(offset, limit, count)
	if start >= 0 && stop >= 0 {
		return data[start:stop]
	}
	return nil
}

func TestListSlice(t *testing.T) {
	t.Log(0, 0, GetElements(0, 0))
	t.Log(0, 1, GetElements(0, 1))
	t.Log(0, -1, GetElements(0, -1))
	t.Log(1, 0, GetElements(1, 0))
	t.Log(1, 1, GetElements(1, 1))
	t.Log(1, -1, GetElements(1, -1))
	t.Log(-1, 0, GetElements(-1, 0))
	t.Log(-1, 1, GetElements(-1, 1))
	t.Log(-1, -1, GetElements(-1, -1))
}
