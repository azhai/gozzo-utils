package geohash

import (
	"math"

	"github.com/azhai/gozzo-utils/random"
	"github.com/kellydunn/golang-geo"
)

const (
	EighthOfCircle = 45        // 八分之一圈的角度
	NauticalMile   = 1.852     // 1海里等于1.852公里
	UnitMps        = 1.0       // 米每秒
	UnitKmph       = 3.6       // 码，公里每小时
	UnitKnot       = 1.9438445 // 节，海里每小时
)

// 通过角度（顺时针，0-360度）获取方位名称
func GetOrientDesc(degree int) string {
	descs := []string{"正北方", "东北方", "正东方",
		"东南方", "正南方", "西南方", "正西方", "西北方"}
	var i int
	sixteenth := EighthOfCircle / 2
	if degree < EighthOfCircle*8-sixteenth {
		i = (degree + sixteenth) / EighthOfCircle
	}
	return descs[i]
}

// 经度或纬度
type Dimension struct {
	Value float64
}

type Position struct {
	Latitude, Longitude float64 // 纬度和经度
	Altitude            float32 // 海拔高度
}

func GetDistance(lastLat, lastLng, currLat, currLng float64) (distance, bearing int) {
	last := geo.NewPoint(lastLat, lastLng)
	curr := geo.NewPoint(currLat, currLng)
	distance = int(last.GreatCircleDistance(curr) * 1000)
	bearing = int(math.Round(last.BearingTo(curr)))
	if bearing < 0 {
		bearing += 360
	}
	return
}

func GetRandSpeed(gap, bearing, angle int, distance float64) (speed float32, orient int) {
	orient = bearing + random.RandMinMax(-15, 15)          // 左右摇摆15度
	ratio := float32(300-angle+random.RandInt(31)) / 200.0 // 根据角度差异，增加5%~70%的速度
	speed = float32(distance) * ratio / float32(gap)
	return
}
