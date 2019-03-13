// Copyright 2018 Ryan Liu. All rights reserved.
// A geohash algorithm of Hilbert space

/*
Example:

import (
	"fmt"
	"geohash"
)

func TestCoordinate() {
	coord := NewCoordinate(10*1000) //距离10km以内
	hash := coord.Encode(22.541497, 113.95196)
	//完整hash为2313000100002333212012
	fmt.Println(hash) //2313000100
}
*/

package geohash

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
)

const (
	LAT_MIN = -90.0
	LAT_MAX = 90.0
	LNG_MIN = -180.0
	LNG_MAX = 180.0
)

//误差表
var precErrors = []float64{
	20015087, //20015.087 km  index=0
	10007543, //10007.543 km
	5003772,  // 5003.772 km
	2501886,  // 2501.886 km
	1250943,  // 1250.943 km
	625471,   //  625.471 km  index=5
	312736,   //  312.736 km
	156368,   //  156.368 km
	78184,    //   78.184 km
	39092,    //   39.092 km
	19546,    //   19.546 km  index=10
	9772.99,  // 9772.992  m
	4886.50,  // 4886.496  m
	2443.25,  // 2443.248  m
	1221.62,  // 1221.624  m
	610.81,   //  610.812  m  index=15
	305.41,   //  305.406  m
	152.70,   //  152.703  m
	76.35,    //   76.351  m
	38.18,    //   38.176  m
	19.09,    //   19.088  m  index=20
	9.54,     //  954.394 cm
	4.77,     //  477.197 cm
}

// 计算前缀长度
func GetCoordPreLen(distance float64) int {
	length := len(precErrors)
	// 二分查找
	i := sort.Search(length, func(i int) bool {
		return precErrors[i] <= distance
	})
	if i < 0 {
		return -1
	} else if i >= length {
		return length
	}
	return i
}

// 经纬坐标点
type Coordinate struct {
	prec uint64
	dim  float64
}

func NewCoordinate(distance float64) *Coordinate {
	const BITS_PER_CHAR = 2
	size := GetCoordPreLen(distance)
	prec := uint64(int64(size))
	dim := 1 << ((prec * BITS_PER_CHAR) >> 1)
	if dim < 1 {
		err := errors.New("Dim must great than or equal 1")
		panic(err)
	}
	return &Coordinate{prec: prec, dim: float64(dim)}
}

func (c *Coordinate) Check(lat, lng float64) bool {
	if lng < LNG_MIN || lng > LNG_MAX {
		return false
	}
	if lat < LAT_MIN || lat > LAT_MAX {
		return false
	}
	return true
}

// Geohash，使用Hilbert空间算法
// prec取值范围1~22，对应误差表中的index
func (c *Coordinate) EncodeInt64(lat, lng float64) int64 {
	if !c.Check(lat, lng) {
		return 0
	}
	x, y := c.coord2Int(lng, lat)
	return c.xy2Hash(int64(x), int64(y))
}

// Geohash，使用Hilbert空间算法
// prec取值范围1~22，对应误差表中的index
func (c *Coordinate) Encode(lat, lng float64) string {
	const (
		_BASE4 = "0123"
		_MASK  = 3
	)
	code := c.EncodeInt64(lat, lng)
	if code <= 0 {
		return ""
	}
	code_size := math.Log2(float64(code)) + 2.0
	code_len := int64(math.Floor(code_size / 2))
	res := make([]byte, code_len)
	for i := code_len - 1; i >= 0; i-- {
		res[i] = _BASE4[code&_MASK]
		code = code >> 2
	}
	length := strconv.FormatUint(c.prec, 10)
	return fmt.Sprintf("%0"+length+"s", string(res))
}

func (c *Coordinate) coord2Int(lng, lat float64) (float64, float64) {
	lngX := (lng + LNG_MAX) / 360.0 * c.dim //[0 ... dim)
	latY := (lat + LAT_MAX) / 180.0 * c.dim //[0 ... dim)
	x := math.Min(c.dim-1, math.Floor(lngX))
	y := math.Min(c.dim-1, math.Floor(latY))
	return x, y
}

func (c *Coordinate) xy2Hash(x, y int64) int64 {
	var (
		rx, ry, d int64
		lvl       = int64(c.dim) >> 1
	)
	for lvl > 0 {
		if (x & lvl) > 0 {
			rx = 1
		} else {
			rx = 0
		}
		if (y & lvl) > 0 {
			ry = 1
		} else {
			ry = 0
		}
		d += lvl * lvl * ((3 * rx) ^ ry)
		x, y = coordRotate(lvl, x, y, rx, ry)
		lvl = lvl >> 1
	}
	return d
}

func coordRotate(n, x, y, rx, ry int64) (int64, int64) {
	if ry == 0 {
		if rx == 1 {
			x = n - 1 - x
			y = n - 1 - y
		}
		return y, x
	}
	return x, y
}
