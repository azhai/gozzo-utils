package calendar

import (
	"time"

	"github.com/azhai/gozzo-utils/common"
)

const (
	LAYOUT_DATE          = "2006-01-02"
	COUNT_WEEK_DAY uint8 = 6 // Weekday分类数量
)

type Weekday = uint8

const (
	W_MON_FRI  Weekday = iota + 1 // 周一到周五
	W_FAKE_SAT                    // 小周六(上班)
	W_HALF_SAT                    // 半天周六
	W_SAT_DAY                     // 大周六(放假)
	W_SUN_DAY                     // 周日
)

type DateKind = uint8

const (
	DK_ILLEGAL  = iota * COUNT_WEEK_DAY // 非法数据
	DK_DAYOFF                           // 被调休日(上班)
	DK_DAY                              // 普通日期
	DK_EVE                              // 副节日
	DK_FESTIVAL                         // 节日
)

// 相差多少天
func GetDiffDays(start, end string) (int, error) {
	starttime, endtime, err := GetTimeRange(start, end)
	if err != nil {
		return 0, err
	}
	secs := endtime.Unix() - starttime.Unix()
	return int(secs / 86400), nil
}

type Calendar struct {
	cal        map[string]uint8
	Start, End time.Time
}

/**
 * 创建一个日历
 * @param saturday_as （第一个）周六等同于哪种类型
 */
func NewCalendar(start, end string, saturday_as Weekday) *Calendar {
	starttime, endtime, err := GetTimeRange(start, end)
	if err != nil {
		return nil
	}
	c := &Calendar{Start: starttime, End: endtime}
	c.Init(saturday_as)
	return c
}

func NewYearCalendar(year int, saturday_as Weekday) *Calendar {
	if year <= 0 {
		year = time.Now().Year()
	}
	c := &Calendar {
		Start: common.NewDate(year, 1, 1),
		End: common.NewDate(year, 12, 31),
	}
	c.Init(saturday_as)
	return c
}

/**
 * 重建日历
 * @param saturday_as （第一个）周六等同于哪种类型
 */
func (c *Calendar) Init(saturday_as Weekday) {
	dt, wd := c.Start, c.Start.Weekday()
	saturday_as = NextSaturday(saturday_as)
	c.cal = make(map[string]uint8) // 清空
	for dt.Before(c.End) {
		date := dt.Format(LAYOUT_DATE)
		if wd == time.Sunday {
			c.cal[date] = W_SUN_DAY + DK_DAY
		} else if wd == time.Saturday {
			saturday_as = NextSaturday(saturday_as)
			c.cal[date] = saturday_as + DK_DAY
		}
		dt = dt.Add(time.Hour * 24)
		wd = NextWeekday(wd)
	}
}

func (c *Calendar) Get(date string) uint8 {
	if val, ok := c.cal[date]; ok {
		return val
	}
	return DK_ILLEGAL
}

func (c *Calendar) SetHoliday(date string) {
	if !c.IsHoliday(date) {
		c.SetDateKind(date, DK_FESTIVAL)
	}
}

func (c *Calendar) SetWorkday(date string) {
	if c.IsHoliday(date) {
		if c.cal[date]%COUNT_WEEK_DAY == W_MON_FRI {
			delete(c.cal, date)
		} else {
			c.SetDateKind(date, DK_DAYOFF)
		}
	}
}

func (c *Calendar) SetDateKind(date string, dk DateKind) uint8 {
	if val, ok := c.cal[date]; ok {
		c.cal[date] = val%COUNT_WEEK_DAY + dk
		return c.cal[date]
	} else if dk > DK_DAY {
		c.cal[date] = W_MON_FRI + dk
		return c.cal[date]
	} else {
		return W_MON_FRI + DK_DAY
	}
}

// 是否放假
func (c *Calendar) IsHoliday(date string) bool {
	return c.Get(date) > W_HALF_SAT+DK_DAY
}

/**
 * 哪些日期放假
 * @param exclude_end 结尾日期不含在内
 */
func (c *Calendar) GetHolidays(start, end string, exclude_end bool) (holidays []string) {
	dt, endtime, err := GetTimeRange(start, end)
	if err != nil {
		return
	}
	if !exclude_end {
		endtime = endtime.Add(time.Hour * 24)
	}
	for dt.Before(endtime) {
		date := dt.Format(LAYOUT_DATE)
		if c.IsHoliday(date) {
			holidays = append(holidays, date)
		}
		dt = dt.Add(time.Hour * 24)
	}
	return
}

func GetTimeRange(start, end string) (starttime, endtime time.Time, err error) {
	starttime, err = time.Parse(LAYOUT_DATE, start)
	if err != nil {
		return
	}
	endtime, err = time.Parse(LAYOUT_DATE, end)
	if err != nil {
		return
	}
	if starttime.After(endtime) {
		starttime, endtime = endtime, starttime
	}
	return
}

func GetWeekday(date string) (time.Weekday, error) {
	dt, err := time.Parse(LAYOUT_DATE, date)
	if err != nil {
		return 0, err
	}
	return dt.Weekday(), nil
}

func NextWeekday(wd time.Weekday) time.Weekday {
	if wd == time.Saturday {
		return 0
	}
	return wd + 1
}

// 大小周切换
func NextSaturday(wd Weekday) Weekday {
	switch wd {
	default:
		return wd
	case W_FAKE_SAT:
		return W_SAT_DAY
	case W_SAT_DAY:
		return W_FAKE_SAT
	}
}
