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
	assert.Equal(t, 2, rh.GetTimeout())
	a, err := rh.GetInt(key)
	assert.NoError(t, err)
	assert.Equal(t, 40, a)
	time.Sleep(2 * time.Second)
	assert.Equal(t, -2, rh.GetTimeout())
}

type Profile struct {
	Age int
	RealName
	Address
}

func NewProfile(id, age int, name RealName, addr Address) *Profile {
	name.ID, addr.ID = id, id
	return &Profile{
		Age: age,
		RealName: name,
		Address: addr,
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

func TestCache(t *testing.T) {
	ryan := NewProfile(5, 40,
		RealName{
			FirstName: "Ryan",
			LastName: "Liu",
		},
		Address{
			City: "深圳",
			Street: "坂田",
		})
	rh := NewRedisHash(GetRedis(), "profile:5", 60)
	ok, err := rh.SetMapJson(Map{
		"age": ryan.Age,
		"name": ryan.RealName,
		"addr": ryan.Address,
	})
	assert.True(t, ok)
	assert.NoError(t, err)
}
