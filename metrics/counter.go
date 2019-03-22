package metrics

import (
	"sync"
	"sync/atomic"
)

// 简单计数器，不能扣除，只能归零
type Counter struct {
	value *uint64
}

func NewCounter() *Counter {
	return &Counter{value: new(uint64)}
}

func (c *Counter) Reset() {
	atomic.StoreUint64(c.value, 0)
}

func (c *Counter) GetCount(name string) uint64 {
	return atomic.LoadUint64(c.value)
}

func (c *Counter) IncrCount(delta uint64) uint64 {
	return atomic.AddUint64(c.value, delta)
}

// 循环列表
type Ring struct {
	count   int
	pointer int
	mutex   *sync.RWMutex
}

func NewRing(count int) *Ring {
	return &Ring{count: count, mutex: new(sync.RWMutex)}
}

func (r *Ring) Next() (curr int) {
	if r.count == 0 {
		return -1
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if r.pointer >= r.count {
		r.pointer = 0
	}
	curr = r.pointer
	r.pointer++
	return
}
