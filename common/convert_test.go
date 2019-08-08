package common

import (
	"encoding/json"
	"testing"
)

type Person struct {
	Name  string
	Birth JsonTime
}

func TestDecimal(t *testing.T) {
	d := NewDecimal(3.1415926, 7)
	dd := ParseDecimal(d.String(), 7)
	diff := d.GetFloat() - dd.GetFloat()
	if diff > 0.0000001 {
		t.Errorf("The diff %.f is too large", diff)
	}
}


func TestJsonTime(t *testing.T) {
	var p Person
	layout := "2006-01-02 15:04:05"
	text := "{\"Name\":\"Ryan\",\"Birth\":\"1981-08-01 12:34:56\"}"
	err := json.Unmarshal([]byte(text), &p)
	if err != nil {
		t.Errorf("Unmarshal fail")
		return
	}
	birth := p.Birth.Format(layout)
	if birth != "1981-08-01 12:34:56" {
		t.Errorf("Unmarshal birth fail %s", birth)
	}
}