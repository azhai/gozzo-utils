package metrics

import (
	"fmt"
	"sync/atomic"
)

type Reporter interface {
	Reset() // 重置、清零
	GetNames() []string
	GetCount(name string) int64
	IncrCount(name string, offset int64) int64
}

//内部计数
type DummyReporter struct {
	names    []string
	counters map[string]*int64
}

func NewDummyReporter(names []string) *DummyReporter {
	r := &DummyReporter{
		names:    names,
		counters: make(map[string]*int64),
	}
	for _, name := range r.names {
		value := int64(0)
		r.counters[name] = &value
	}
	return r
}

func (r *DummyReporter) Reset() {
	for _, value := range r.counters {
		atomic.StoreInt64(value, 0)
	}
}

func (r *DummyReporter) GetNames() []string {
	return r.names
}

func (r *DummyReporter) GetCount(name string) int64 {
	if value, ok := r.counters[name]; ok {
		return atomic.LoadInt64(value)
	}
	return 0
}

func (r *DummyReporter) IncrCount(name string, offset int64) int64 {
	if value, ok := r.counters[name]; ok {
		return atomic.AddInt64(value, offset)
	}
	return 0
}

//输出当前统计结果
func StatSnap(r Reporter, sameWidth bool) string {
	var result string = ""
	tpl := " %s=%d"
	if sameWidth {
		tpl = "  %s=% 8d"
	}
	names := r.GetNames()
	for _, name := range names {
		if value := r.GetCount(name); value >= 0 {
			result += fmt.Sprintf(tpl, name, value)
		}
	}
	return result
}
