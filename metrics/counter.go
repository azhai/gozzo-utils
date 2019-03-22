package metrics

import "sync/atomic"

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
