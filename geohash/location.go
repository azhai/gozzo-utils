package geohash

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/random"
	"github.com/kellydunn/golang-geo"
)

const (
	EighthOfCircle = 45                 // 八分之一圈的角度
	NauticalMile   = float32(1.852)     // 1海里等于1.852公里
	// 速度单位
	UnitMps        = float32(1.0)       // 米每秒
	UnitKmph       = float32(3.6)       // 码，公里每小时
	UnitKnot       = float32(1.9438445) // 节，海里每小时
	// 投影点位置判断
	OutsideB       = -2 // B点外侧
	OutsideA       = -1 // A点外侧
	CoincideA      = 0  // 与A点重合
	InsideA        = 1  // 内侧靠近A点（含中间点）
	InsideB        = 2  // 内侧靠近B点
	CoincideB      = 3  // 与B点重合
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

// 平面投影：将C点按比例投影到AB直线上，得到点D以及更靠近哪一个点
func MapPoint(a, b, c *geo.Point) (*geo.Point, int) {
	dLat, dLng := b.Lat() - a.Lat(), b.Lng() - a.Lng()
	diag := math.Sqrt(dLat * dLat + dLng * dLng) // 对角线
	dx, dy := dLat / diag, dLng / diag
	dist := dx * (c.Lat() - a.Lat()) + dy * (c.Lng() - a.Lng())
	lat, lng := a.Lat() + dist * dx, a.Lng() + dist * dy
	d := geo.NewPoint(common.RoundN(lat, 6), common.RoundN(lng, 6)) // 投影
	delta := math.Pow10(-6)
	if dist < 0.0 - delta {
		return d, OutsideA
	} else if dist <= delta {
		return d, CoincideA
	}
	diffWhole := dist - float64(diag)
	if diffWhole > delta {
		return d, OutsideB
	} else if diffWhole >= 0.0 - delta {
		return d, CoincideB
	}
	diffHalf := dist - float64(diag) / 2.0
	if diffHalf <= delta {
		return d, InsideA
	} else {
		return d, InsideB
	}
}

// 围栏接口
type Fence interface {
	Contains(point *geo.Point) bool // 是否在围栏内（含边界）
}

// 圆形围栏
type Circle struct {
	Center *geo.Point // 中心点
	Radius int        // 半径
}

func (c *Circle) Contains(point *geo.Point) bool {
	if c.Radius <= 0 {
		return false
	}
	distKilo := c.Center.GreatCircleDistance(point)
	return int(distKilo*1000) <= c.Radius
}

// 多边形围栏
type Polygon = geo.Polygon

// 航线围栏
type Stripe struct {
	padding int
	points map[string]*geo.Point
	values []string
	coord  *Coordinate
}

// padding为道路单边宽度（米）
func NewStripe(padding int, points []*geo.Point) *Stripe {
	s := &Stripe{padding, make(map[string]*geo.Point), nil, nil}
	s.Insert(points...)
	return s
}

// 计算Hilbert哈希值
func (s *Stripe) Hash(point *geo.Point) string {
	if s.coord == nil {
		s.coord = NewCoordinate(float64(s.padding))
	}
	return s.coord.Encode(point.Lat(), point.Lng())
}

func (s *Stripe) Len() int {
	return len(s.values)
}

func (s *Stripe) Values() []string {
	return s.values
}

func (s *Stripe) Point(i int) *geo.Point {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return s.points[s.values[i]]
}

// 在合适位置增加多个点
func (s *Stripe) Insert(points ...*geo.Point) {
	// 将原有值和新增的值放在一起排序，重新构建列表
	for _, p := range points {
		newbie := s.Hash(p)
		s.values = append(s.values, newbie)
		s.points[newbie] = p
	}
	sort.Strings(s.values)
}

// 找出最近的两个点
func (s *Stripe) Nearest(point *geo.Point) (a, b *geo.Point) {
	value := s.Hash(point)
	// 二分查找，从前向后查找，idx可能是0-len
	idx := sort.Search(s.Len(), func(i int) bool {
		return strings.Compare(s.values[i], value) >= 0
	})
	if idx == 0 { // 可能小于或等于首个元素
		b = s.Point(0)
	} else if idx == s.Len() {
		a = s.Point(idx - 1)
	} else {
		a, b = s.Point(idx - 1), s.Point(idx)
	}
	return
}

func (s *Stripe) Contains(point *geo.Point) bool {
	var shadow *geo.Point
	a, b := s.Nearest(point)
	if a == nil {
		shadow = b
	} else if b == nil {
		shadow = a
	} else {
		var near int
		shadow, near = MapPoint(a, b, point)
		if near == CoincideA || near == CoincideB {
			return true
		}
	}
	dist, _ := GetDistance(shadow, point)
	return dist <= s.padding
}
