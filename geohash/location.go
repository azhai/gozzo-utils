package geohash

import (
	"math"

	"github.com/azhai/gozzo-utils/common"
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

// 计算A点到B点的距离（单位：米）和方向（顺时针0-359度）
func GetDistance(latA, lngA, latB, lngB float64) (int, int) {
	posA, posB := geo.NewPoint(latA, lngA), geo.NewPoint(latB, lngB)
	distance := int(posA.GreatCircleDistance(posB) * 1000)
	bearing := int(math.Round(posA.BearingTo(posB)))
	if bearing < 0 {
		bearing += 360
	}
	return distance, bearing
}

// 根据角度差和时间（秒），估算行驶这段距离的速度（公里每小时）和方向
// bearing为当前方向（比如B到C），oldBearing为旧的方向（比如A到B）
func GetInexactSpeed(distance, bearing, oldBearing, gap int) (float32, int) {
	angle := (90 - (bearing-oldBearing)%180) % 90 // 0-90的角度差，相差越大，实际距离比直线距离越大
	ratio := float32(300-angle+random.RandInt(31)) / 200.0 // 根据角度差异，增加5%~70%的速度
	speed := float32(distance) * ratio * UnitKmph / float32(gap)
	bearing += random.RandMinMax(-15, 15) // 左右摇摆15度
	return speed, bearing
}

// 经度或纬度
type Dimension struct {
	*common.Decimal
}

func NewDimension(value float64) *Dimension {
	return &Dimension{common.NewDecimal(value, 6)}
}

type Position struct {
	Latitude, Longitude *Dimension // 纬度和经度
	Altitude            float32 // 海拔高度
}
