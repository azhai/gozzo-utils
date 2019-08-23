package geohash

import (
	"testing"

	"github.com/kellydunn/golang-geo"
	"github.com/stretchr/testify/assert"
)

/*
深圳市莲花山公园
打开网站 http://www.gpsspg.com/distance.htm
先点击地图上方的第二个按钮“清空第题（不记录）”
然后在左侧输入框中粘贴下面的数据

22.5478429411,114.0443715561
22.5494580855,114.0443179119
22.5567557264,114.0439048517
22.5577663672,114.0433791387
22.5599988198,114.0477065541
22.5618585457,114.0500122771

最后点击下方的“提交”按钮，可看到一条沿着莲花山西路线
*/

var (
	radius = 1130                         // 米
	center = "22.556766583,114.053219182" // 中心
	outer  = "22.5512664303,114.0436741817"
	inner  = "22.5567557264,114.0439048517" // 正西
	road   = []string{
		"22.5478429411,114.0443715561",
		"22.5494580855,114.0443179119",
		"22.5567557264,114.0439048517", // 正西
		"22.5577663672,114.0433791387",
		"22.5599988198,114.0477065541",
		"22.5618585457,114.0500122771", // 西北角
	}
	points = []string{
		"22.5623886180,114.0579891667", // 东北角
		"22.5598491952,114.0605899327",
		"22.5583827955,114.0618237488",
		"22.5571361918,114.0625335120", // 正东
		"22.5547482831,114.0625656985",
		"22.5532025609,114.0623511218",
		"22.5515280089,114.0625549697", // 东南角
		"22.5513823517,114.0557485099",
		"22.5513625344,114.0507810589",
		"22.5510781654,114.0477404106",
		"22.5509097176,114.0458467710",
		"22.5507610871,114.0443500984", // 西南角
		"22.5538971581,114.0442857254",
		"22.5556063612,114.0442159879",
		"22.5567557264,114.0439048517", // 正西
		"22.5599988198,114.0477065541",
		"22.5618585457,114.0500122771", // 西北角
		"22.5623985259,114.0511012539",
		"22.5625768674,114.0517718062",
		"22.5623886180,114.0541965231", // 正北
	}
)

func TestGeoHash(t *testing.T) {
	coord := NewCoordinate(5) //距离5m以内
	point := ToPoint(center)
	hash := coord.Encode(point.Lat(), point.Lng())
	assert.Equal(t, "2313000100023130213230", hash)
	lat, lng := coord.Decode(hash)
	assert.InDelta(t, lat, point.Lat(), 2e-05)
	assert.InDelta(t, lng, point.Lng(), 2e-05)
	point2 := geo.NewPoint(lat, lng)
	dist, _ := GetDistance(point, point2)
	assert.True(t, dist <= 5)
	t.Log(point, hash, dist)
}

func TestGetDistance(t *testing.T) {
	a, b := ToPoint(points[0]), ToPoint(points[1])
	dist, bear := GetDistance(a, b)
	assert.Equal(t, 388, dist)
	assert.Equal(t, 137, bear)
}

func TestMapPoint(t *testing.T) {
	a, b := ToPoint(points[2]), ToPoint(points[3])
	c, d := ToPoint(points[4]), ToPoint(points[5])
	t.Log("a: ", a)
	t.Log("b: ", b)
	t.Log("c: ", c)
	t.Log("d: ", d)
	shadow, between := MapPoint(a, c, a)
	assert.Equal(t, between, CoincideA)
	t.Log("aca", shadow, between)
	shadow, between = MapPoint(a, c, b)
	assert.Equal(t, between, InsideA)
	t.Log("acb", shadow, between)
	shadow, between = MapPoint(a, c, d)
	assert.Equal(t, between, OutsideB)
	t.Log("acd", shadow, between)
}

func TestCircleFence(t *testing.T) {
	f := &Circle{Center: ToPoint(center), Radius: radius}
	assert.True(t, f.Contains(ToPoint(inner)))
	assert.False(t, f.Contains(ToPoint(outer)))
}

func TestPolygonFence(t *testing.T) {
	var ps []*geo.Point
	for _, p := range road {
		ps = append(ps, ToPoint(p))
	}
	f := geo.NewPolygon(ps)
	assert.True(t, f.Contains(ToPoint(inner)))
	assert.False(t, f.Contains(ToPoint(outer)))
}

func TestStripeFence(t *testing.T) {
	var ps []*geo.Point
	for _, p := range road {
		ps = append(ps, ToPoint(p))
	}
	f := NewStripe(20, ps)
	for i := 0; i < f.Len(); i++ {
		t.Log("pn: ", i, f.Point(i))
	}

	shadow, between := f.NearestPoints(ToPoint(inner))
	t.Log("inner", ToPoint(inner))
	t.Log(shadow, between)
	assert.True(t, f.Contains(ToPoint(inner)))

	shadow, between = f.NearestPoints(ToPoint(outer))
	t.Log("outer", ToPoint(outer))
	t.Log(shadow, between)
	assert.False(t, f.Contains(ToPoint(outer)))
}
