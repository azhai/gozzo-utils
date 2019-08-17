package geohash

import (
	"testing"

	"github.com/azhai/gozzo-utils/common"
	"github.com/kellydunn/golang-geo"
	"github.com/stretchr/testify/assert"
)

// 深圳市莲花山公园
var (
	radius = 1130                         // 米
	center = "22.556766583,114.053219182" // 中心
	outer  = "22.5512664303,114.0436741817"
	inner  = "22.5567557264,114.0439048517" // 正西
	road   = []string{
		"22.5478429411,114.0443715561",
		"22.5494580855,114.0443179119",
		"22.5567557264,114.0439048517", // 正西
		"22.5599988198,114.0477065541",
		"22.5618585457,114.0500122771", // 西北角
		"22.5577663672,114.0433791387",
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

func ToPoint(p string) *geo.Point {
	pieces := common.SplitPieces(p, ",", nil)
	if len(pieces) < 2 {
		return nil
	}
	lat := common.ParseDecimal(pieces[0], 6)
	lng := common.ParseDecimal(pieces[1], 6)
	return geo.NewPoint(lat.GetFloat(), lng.GetFloat())
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
	t.Log(f.Hash(ToPoint(outer))[11:])
	t.Log(f.Find(ToPoint(outer)))
	t.Log(f.Point(4))
	t.Log(f.Point(5))
	for _, v := range f.Values() {
		t.Log(v[11:])
	}
	assert.True(t, f.Contains(ToPoint(inner)))
	assert.False(t, f.Contains(ToPoint(outer)))
}
