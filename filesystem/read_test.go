package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	fname := "./read.go"
	count := LineCount(fname)
	t.Log(count, fname)

	// 逐行返回，适用于大文件
	var lines []string
	r := NewLineReader(fname)
	for r.Reading() {
		lines = append(lines, r.Text())
	}
	err := r.Err()
	assert.NoError(t, err)
	assert.Len(t, lines, count)

	// 直接返回全部行，适用于小文件
	lines, err = ReadLines(fname)
	assert.NoError(t, err)
	assert.Len(t, lines, count)
}

