package calendar

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func SetCalendarY2019(cal *Calendar) *Calendar {
	// 元旦
	cal.SetHoliday("2019-01-01")
	// 春节
	cal.SetWorkday("2019-02-02")
	cal.SetWorkday("2019-02-03")
	cal.SetHoliday("2019-02-04")
	cal.SetHoliday("2019-02-05")
	cal.SetHoliday("2019-02-06")
	cal.SetHoliday("2019-02-07")
	cal.SetHoliday("2019-02-08")
	// 清明节
	cal.SetHoliday("2019-04-05")
	// 劳动节
	cal.SetWorkday("2019-04-28")
	cal.SetHoliday("2019-05-01")
	cal.SetHoliday("2019-05-02")
	cal.SetHoliday("2019-05-03")
	cal.SetWorkday("2019-05-05")
	// 端午节
	cal.SetHoliday("2019-06-07")
	// 中秋节
	cal.SetHoliday("2019-09-13")
	// 国庆节
	cal.SetWorkday("2019-09-29")
	cal.SetHoliday("2019-10-01")
	cal.SetHoliday("2019-10-02")
	cal.SetHoliday("2019-10-03")
	cal.SetHoliday("2019-10-04")
	cal.SetHoliday("2019-10-07")
	cal.SetWorkday("2019-10-12")
	return cal
}

func TestCalcDay(t *testing.T) {
	assert.False(t, IsSolarLeap(1900))
	assert.False(t, IsSolarLeap(1990))
	assert.False(t, IsSolarLeap(1994))
	assert.True(t, IsSolarLeap(2004))
	diff, err := GetDiffDays("2019-01-05", "2019-05-01")
	assert.NoError(t, err)
	assert.Equal(t, 116, diff)
	diff, err = GetDiffDays("2019-01-05", "2019-10-05")
	assert.NoError(t, err)
	assert.Equal(t, 273, diff)
}

func TestMidAutumn(t *testing.T) {
	autumn := NewLunar(2016, 8, 15)
	for i := 2016; i <= 2020; i++ {
		autumn.LunarYear = i
		dt := FormatSolar(LunarToSolar(autumn))
		t.Logf("MidAutumn of year %d: %s", i, dt)
		assert.Equal(t, MidAutumn.GetFirstDate(i), dt)
	}
}

func TestGetHolidays(t *testing.T) {
	cal := NewYearCalendar(2019, W_FAKE_SAT)
	assert.Equal(t, cal.Start.Month(), time.January)
	assert.Equal(t, cal.Start.Day(), 1)
	assert.Equal(t, cal.End.Month(), time.December)
	assert.Equal(t, cal.End.Day(), 31)
	holidays := cal.GetHolidays("2019-05-01", "2019-05-31", false)
	t.Logf("Before setting: %+v", holidays)
	assert.Len(t, holidays, 6)
	cal = SetCalendarY2019(cal)
	holidays = cal.GetHolidays("2019-05-01", "2019-05-31", false)
	t.Logf("After setting: %+v", holidays)
	assert.Len(t, holidays, 8)
}

func TestFakeSaturday(t *testing.T) {
	cal := NewYearCalendar(2019, W_FAKE_SAT)
	cal = SetCalendarY2019(cal)
	fifthes := []bool{
		false, true, false, // 小周六、周二&春节初二、周二
		true, false, false, // 周五&清明节、周日&被调休、周三
		false, false, false, // 周五、周一、周四
		true, false, false, // 大周六&国庆节、周二、周四
	}
	for i := 1; i <= 12; i++ {
		dt := fmt.Sprintf("2019-%02d-05", i)
		t.Logf("%s: %v", dt, cal.IsHoliday(dt))
		assert.Equal(t, cal.IsHoliday(dt), fifthes[i-1])
	}
}
