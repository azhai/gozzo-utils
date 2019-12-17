package calendar

import (
	"fmt"
	"time"
)

var (
	NewYear   = NewSolarFestival(1, 1, "元旦")
	// 春节放假的第一天实际上是除夕
	SpringDay = NewLunarFestivalCount("春节", 3).AddAnnals(2016,
		"2016-02-07", "2017-01-27", "2018-02-15", "2019-02-04", "2020-01-24")
	// 清明节有时是4月4日，放假的第一天是4月2日，例如2016年
	TombSweeping = NewSolarFestival(4, 5, "清明节")
	LabourDay    = NewSolarFestival(5, 1, "劳动节")
	DragonBoat   = NewLunarFestival("端午节").AddAnnals(2016,
		"2016-06-09", "2017-05-28", "2018-06-18", "2019-06-07", "2020-06-25")
	MidAutumn = NewLunarFestival("中秋节").AddAnnals(2016,
		"2016-09-15", "2017-10-04", "2018-09-24", "2019-09-13", "2020-10-01")
	NationalDay = NewSolarFestival(10, 1, "国庆节").SetCount(3)
)

func NewLunar(year, month, day int) Lunar {
	return Lunar{
		LunarYear:  year,
		LunarMonth: month,
		LunarDay:   day,
	}
}

func NewSolar(date string) Solar {
	dt, err := time.Parse(LAYOUT_DATE, date)
	if err != nil {
		return Solar{}
	}
	return Solar{
		SolarYear:  dt.Year(),
		SolarMonth: int(dt.Month()),
		SolarDay:   dt.Day(),
	}
}

func FormatSolar(s *Solar) string {
	return fmt.Sprintf("%04d-%02d-%02d", s.SolarYear, s.SolarMonth, s.SolarDay)
}

func IsSolarLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// 节日
type FestivalDay struct {
	Title    string
	DayCount int
}

func NewFestivalDay(title string) FestivalDay {
	return FestivalDay{Title: title, DayCount: 1}
}

func (f FestivalDay) GetTitle() string {
	return f.Title
}

func (f FestivalDay) GetDays() int {
	return f.DayCount
}

// 法定节日
type Festival interface {
	GetTitle() string
	GetDays() int
	GetFirstDate(year int) string
}

// 公历节日
type SolarFestival struct {
	Month    int
	FirstDay int
	FestivalDay
}

func NewSolarFestival(month, day int, title string) SolarFestival {
	return SolarFestival{
		Month: month, FirstDay: day,
		FestivalDay: NewFestivalDay(title),
	}
}

func (f SolarFestival) SetCount(days int) SolarFestival {
	f.DayCount = days
	return f
}

func (f SolarFestival) GetFirstDate(year int) string {
	return fmt.Sprintf("%04d-%02d-%02d", year, f.Month, f.FirstDay)
}

// 农历节日
type LunarFestival struct {
	Annals map[int]string // 年鉴
	FestivalDay
}

func NewLunarFestival(title string) LunarFestival {
	return LunarFestival{
		Annals:      make(map[int]string),
		FestivalDay: NewFestivalDay(title),
	}
}

func NewLunarFestivalCount(title string, days int) LunarFestival {
	f := NewLunarFestival(title)
	f.DayCount = days
	return f
}

func (f LunarFestival) AddAnnals(year int, firstDates ...string) LunarFestival {
	for i, dt := range firstDates {
		f.Annals[year+i] = dt
	}
	return f
}

func (f LunarFestival) GetFirstDate(year int) string {
	if dt, ok := f.Annals[year]; ok {
		return dt
	}
	return ""
}
