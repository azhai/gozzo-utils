package geohash

import (
	"math"
	"time"

	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/random"
	"github.com/kellydunn/golang-geo"
)

const (
	EighthOfCircle = 45             // 八分之一圈的角度
	NauticalMile   = float32(1.852) // 1海里等于1.852公里
	// 速度单位
	UnitMps  = float32(1.0)       // 米每秒
	UnitKmph = float32(3.6)       // 码，公里每小时
	UnitKnot = float32(1.9438445) // 节，海里每小时
	// 投影点位置判断
	OutsideB  = -2 // B点外侧
	OutsideA  = -1 // A点外侧
	CoincideA = 0  // 与A点重合
	InsideA   = 1  // 内侧靠近A点（含中间点）
	InsideB   = 2  // 内侧靠近B点
	CoincideB = 3  // 与B点重合
)

// 通过角度（顺时针，0-360度）获取方位名称
func GetBearingDesc(bearing int) string {
	descs := []string{"正北方", "东北方", "正东方",
		"东南方", "正南方", "西南方", "正西方", "西北方"}
	var i int
	sixteenth := EighthOfCircle / 2
	if bearing < EighthOfCircle*8-sixteenth {
		i = (bearing + sixteenth) / EighthOfCircle
	}
	return descs[i]
}

// 判断线段AB上点C的相对位置，dist是AC距离，whole是AB长度
func GetPointSide(dist, whole, delta float64) int {
	if dist < 0.0-delta {
		return OutsideA
	} else if dist <= delta {
		return CoincideA
	}
	if whole < 0 { // 没有B点，结束判断
		return OutsideA
	}
	// 继续判断与B的位置关系
	diffWhole := dist - float64(whole)
	if diffWhole > delta {
		return OutsideB
	} else if diffWhole >= 0.0-delta {
		return CoincideB
	}
	diffHalf := dist - float64(whole)/2.0
	if diffHalf <= delta {
		return InsideA
	} else {
		return InsideB
	}
}

// 速度换算
type Speed struct {
	Mps, Kmph, Knot float32
}

func NewSpeed(value float32) Speed {
	return Speed{
		Mps:  value * UnitMps,
		Kmph: value * UnitKmph,
		Knot: value * UnitKnot,
	}
}

type Position struct {
	Moment            int64 // 时间戳
	Altitude, Bearing int   // 海拔高度（米）和方向（顺时针0-359度）
	*geo.Point
}

func NewPosition(lat, lng float64, t *time.Time) *Position {
	return &Position{t.Unix(), 0, 0, geo.NewPoint(lat, lng)}
}

// 根据距离和夹角计算另一个坐标点
func (p *Position) Add(dist, angle, alt int) *Position {
	return &Position{
		p.Moment, p.Altitude + alt, 0,
		p.PointAtDistanceAndBearing(float64(dist)/1000.0, float64(angle)),
	}
}

// 计算A点到B点的直线距离（米）和夹角
func (p *Position) GetDistance(target *Position) (int, int) {
	return GetDistance(p.Point, target.Point)
}

// 根据夹角和时间（秒），估算行驶这段距离的速度（米每秒）和方向
func (p *Position) GetInexactSpeed(target *Position) (float32, int) {
	bear := target.Bearing + random.RandMinMax(-15, 15) // 左右摇摆15度
	if p.Moment == target.Moment {
		return 0.0, bear
	}
	if p.Moment > target.Moment {
		return target.GetInexactSpeed(p)
	}
	dist, angle := p.GetDistance(target)
	angle = (90 - (angle-p.Bearing)%180) % 90              // 0-90的角度差，相差越大，实际距离比直线距离越大
	ratio := float32(300-angle+random.RandInt(31)) / 200.0 // 根据角度差异，增加5%~70%的速度
	gap := target.Moment - p.Moment
	speed := float32(dist) * ratio / float32(gap)
	return speed, bear
}

// 角度转弧度
func Bear2Arc(bear int) float64 {
	return float64(bear) * math.Pi / 180.0
}

// 将字符串表示lat,lng的转为Point
func ToPoint(p string) *geo.Point {
	pieces := common.SplitPieces(p, ",", nil)
	if len(pieces) < 2 {
		return nil
	}
	lat := common.ParseDecimal(pieces[0], 6)
	lng := common.ParseDecimal(pieces[1], 6)
	return geo.NewPoint(lat.GetFloat(), lng.GetFloat())
}

// 计算距离：获取距离（米）和角度（顺时针）
func GetDistance(a, b *geo.Point) (int, int) {
	dist := int(a.GreatCircleDistance(b) * 1000)
	angle := int(math.Round(a.BearingTo(b)))
	if angle < 0 {
		angle += 360
	}
	return dist, angle
}
