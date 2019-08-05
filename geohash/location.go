package geohash

import (
	"encoding/binary"
	"math"
	"sort"
	"strings"
	"time"
	"container/list"

	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/random"
	"github.com/kellydunn/golang-geo"
)

const (
	EighthOfCircle = 45                 // 八分之一圈的角度
	NauticalMile   = float32(1.852)     // 1海里等于1.852公里
	UnitMps        = float32(1.0)       // 米每秒
	UnitKmph       = float32(3.6)       // 码，公里每小时
	UnitKnot       = float32(1.9438445) // 节，海里每小时
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

// 经度或纬度
type Dimension struct {
	*common.Decimal
}

func NewDimension(value float64) *Dimension {
	return &Dimension{common.NewDecimal(value, 6)}
}

func (d Dimension) Encode(v interface{}) []byte {
	dim := v.(*Dimension)
	chunk := make([]byte, 4)
	binary.BigEndian.PutUint32(chunk, uint32(dim.Value))
	return chunk
}

func (d Dimension) Decode(chunk []byte) interface{} {
	dim := NewDimension(0.0)
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 4)
		value := binary.BigEndian.Uint32(chunk)
		dim.Value = int64(value)
	}
	return dim
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
	dist := int(p.GreatCircleDistance(target.Point) * 1000)
	angle := int(math.Round(p.BearingTo(target.Point)))
	if angle < 0 {
		angle += 360
	}
	return dist, angle
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

// 围栏接口
type Enclosure interface {
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
	coord *Coordinate
	values []string
	*list.List
}

// padding为道路单边宽度（米）
func NewStripe(padding int, points ...*geo.Point) *Stripe {
	coord := NewCoordinate(float64(padding))
	s := &Stripe{coord, nil, list.New()}
	s.Insert(points...)
	return s
}

// 计算Hilbert哈希值
func (s *Stripe) Hash(point *geo.Point) string {
	return s.coord.Encode(point.Lat(), point.Lng())
}

func (s *Stripe) Values() []string {
	return s.values
}

// 从后往前查找最近的两个点
func (s *Stripe) Seek(value string) (*list.Element, *list.Element) {
	for mark := s.Back(); mark != nil; mark = mark.Prev()  {
		flag := strings.Compare(value, mark.Value.(string))
		if flag == 0 { // 找到一个重复的点
			return nil, mark // 前一个元素留空，方便后续判断
		}
		if flag > 0 { // 与终点同一侧，继续
			continue
		}
		return mark, mark.Next()
	}
	return nil, nil
}

// 在合适位置增加多个点
func (s *Stripe) Insert(points ...*geo.Point) {
	// 将原有值和新增的值放在一起排序，重新构建列表
	for _, p := range points {
		s.values = append(s.values, s.Hash(p))
	}
	sort.Strings(s.values)
	s.Init() // 清空列表
	for _, value := range s.values {
		s.PushBack(value)
	}
}

func (s *Stripe) Contains(point *geo.Point) bool {
	value := s.Hash(point)
	// 二分查找，从前向后查找，idx可能是0-len
	idx := sort.Search(s.Len(), func(i int) bool {
		return strings.Compare(s.values[i], value) >= 0
	})
	if idx == 0 { // 可能小于或等于首个元素
		return strings.Compare(s.values[idx], value) == 0
	}
	return idx < s.Len()
}
