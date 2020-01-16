package redisw

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetRedis() *RedisWrapper {
	return NewRedisPool(ConnParams{}, 0)
}

func TestInt(t *testing.T) {
	name := "test:a"
	r := GetRedis()
	r.SetVal(name, 39, 60)
	assert.Equal(t, 60, r.GetTimeout(name))
	a, err := r.GetInt(name)
	assert.NoError(t, err)
	assert.Equal(t, 39, a)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 58, r.GetTimeout(name))
}

func TestHash(t *testing.T) {
	name, key := "test:hash", "a"
	rh := NewRedisHash(GetRedis(), name, 2)
	rh.SetVal(key, 40)
	assert.Equal(t, 2, rh.GetTimeout(false))
	a, err := rh.GetInt(key)
	assert.NoError(t, err)
	assert.Equal(t, 40, a)
	time.Sleep(2 * time.Second)
	assert.Equal(t, -2, rh.GetTimeout(false))
}

type Profile struct {
	Age int
	*RealName
	*Address
}

func NewProfile(id, age int, name RealName, addr Address) *Profile {
	name.ID, addr.ID = id, id
	return &Profile{
		Age: age,
		RealName: &name,
		Address: &addr,
	}
}

type RealName struct {
	ID int `json:"-"`
	FirstName string `json:"first"`
	LastName string `json:"last"`
}

func (n RealName) GetCacheId() string {
	return fmt.Sprintf("name:%d", n.ID)
}

type Address struct {
	ID int `json:"-"`
	Province string `json:"province"`
	City string `json:"city"`
	Street string `json:"street"`
	Building string `json:"building"`
	Room string `json:"room"`
}

func (a Address) GetCacheId() string {
	return fmt.Sprintf("addr:%d", a.ID)
}

var ryan = NewProfile(5, 40,
	RealName{
		FirstName: "Ryan",
		LastName: "Liu",
	},
	Address{
		City: "深圳",
		Street: "坂田",
	})

func TestJson(t *testing.T) {
	key := "test:name"
	r := GetRedis()
	ok, err := r.SaveJson(key, ryan.RealName, 60)
	assert.True(t, ok)
	assert.NoError(t, err)
	var name = new(RealName)
	t.Logf("before: %#v\n", name)
	err = r.LoadJson(key, name)
	t.Logf("after: %#v\n", name)
	assert.NoError(t, err)
	assert.Equal(t, "Ryan", name.FirstName)
	assert.Equal(t, "Liu", name.LastName)
}

func TestCache(t *testing.T) {
	data := Map{
		"age": ryan.Age,
		"name": ryan.RealName,
		"addr": ryan.Address,
	}
	rh := NewRedisHash(GetRedis(), "profile:5", 60)
	ok, err := rh.SaveMapData(data)
	assert.True(t, ok)
	assert.NoError(t, err)
	data["age"] = 14
	data["addr"] = &Address{
		City: "东莞",
		Street: "松山湖",
	}
	t.Logf("before: %#v\n", data)
	err = rh.LoadMapJson(data)
	var others map[string]int
	others, err = rh.LoadMapInt("age")
	for key, val := range others {
		data[key] = val
	}
	t.Logf("after: %#v\n", data)
	assert.NoError(t, err)
	assert.Equal(t, data["age"], ryan.Age)
	var street string
	if addr, ok := data["addr"].(*Address); ok {
		street = addr.Street
	}
	assert.Equal(t, street, ryan.Address.Street)
}
