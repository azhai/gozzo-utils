package geohash

import (
	"math"
	"sort"
	"strings"

	"github.com/azhai/gozzo-utils/common"
	"github.com/kellydunn/golang-geo"
)

// 使用正弦函数计算投影距离（米）
func GetMapDistance(a, b, c *geo.Point) int {
	_, bearB := GetDistance(a, b)
	distC, bearC := GetDistance(a, c)
	arc := math.Abs(Bear2Arc(bearB - bearC)) // 弧度
	return int(math.Sin(arc) * float64(distC))
}

// 平面投影：将C点按比例投影到AB直线上，得到点D以及更靠近A或B哪一个
func MapPoint(a, b, c *geo.Point) (*geo.Point, int) {
	dLat, dLng := b.Lat()-a.Lat(), b.Lng()-a.Lng()
	diag := math.Sqrt(dLat*dLat + dLng*dLng) // 对角线长度，即AB距离
	dx, dy := dLat/diag, dLng/diag
	dist := dx*(c.Lat()-a.Lat()) + dy*(c.Lng()-a.Lng()) // 投影距离，即AD距离
	lat, lng := a.Lat()+dist*dx, a.Lng()+dist*dy
	d := geo.NewPoint(common.RoundN(lat, 6), common.RoundN(lng, 6)) // 投影点D
	return d, GetPointSide(dist, diag, 1e-6)
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
	padding int // 道路单边宽度（米）
	points  map[string]*geo.Point
	values  []string
	coord   *Coordinate
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

func (s *Stripe) Nearest(value string) (int, int) {
	// 二分查找，从前向后查找，idx可能是0-len
	idx := sort.Search(s.Len(), func(i int) bool {
		return strings.Compare(s.values[i], value) >= 0
	})
	if idx == 0 { // 可能小于或等于首个元素
		return 0, -1
	} else if idx == s.Len() {
		return idx - 1, -1
	} else {
		return idx, idx - 1
	}
}

// 找出最近的两个点
func (s *Stripe) NearestPoints(point *geo.Point) (a, b *geo.Point) {
	value := s.Hash(point)
	i, j := s.Nearest(value)
	a = s.Point(i)
	if j > 0 {
		b = s.Point(j)
	}
	return
}

func (s *Stripe) Contains(point *geo.Point) bool {
	var dist int
	a, b := s.NearestPoints(point)
	if b == nil {
		dist, _ = GetDistance(a, point)
	} else {
		dist = GetMapDistance(a, b, point)
	}
	return dist <= s.padding
}
