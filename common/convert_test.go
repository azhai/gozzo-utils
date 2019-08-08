package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string
	Birth JsonTime
}

func TestDecimal(t *testing.T) {
	d := NewDecimal(3.1415926, 7)
	dd := ParseDecimal(d.String(), 7)
	assert.Zero(t, d.GetFloat()-dd.GetFloat())
}

func TestJsonTime(t *testing.T) {
	var p Person
	layout := "2006-01-02 15:04:05"
	text := "{\"Name\":\"Ryan\",\"Birth\":\"1981-08-01 12:34:56\"}"
	err := json.Unmarshal([]byte(text), &p)
	assert.NoError(t, err, "Unmarshal fail")
	if err == nil {
		birth := p.Birth.Format(layout)
		assert.Equal(t, birth, "1981-08-01 12:34:56")
	}
}
