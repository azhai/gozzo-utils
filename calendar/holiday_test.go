package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMidAutumn(t *testing.T) {
	autumn := NewLunar(2016, 8, 15)
	for i := 2016; i <= 2020; i++ {
		autumn.LunarYear = i
		dt := FormatSolar(LunarToSolar(autumn))
		t.Logf("MidAutumn of year %d: %s", i, dt)
		assert.Equal(t, MidAutumn.GetFirstDate(i), dt)
	}
}

func TestFakeSaturday(t *testing.T) {
	cal := NewYearCalendar(W_FAKE_SAT)
	assert.Equal(t, cal.Start.Month(), time.January)
	assert.Equal(t, cal.Start.Day(), 1)
	assert.Equal(t, cal.End.Month(), time.December)
	assert.Equal(t, cal.End.Day(), 31)
}
