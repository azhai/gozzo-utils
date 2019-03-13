package metrics

import (
	"testing"
)

func CreateReporter() Reporter {
	var names = []string{"dial", "send", "recv"}
	return NewDummyReporter(names)
}

// 测试是否有覆盖的情况
func BenchmarkDummy(b *testing.B) {
	reporter := CreateReporter()
	for i := 0; i < b.N; i++ {
		go func(reporter Reporter) {
			reporter.IncrCount("dial", 1)
			reporter.IncrCount("send", 2)
			reporter.IncrCount("recv", 2)
		}(reporter)
	}
	b.Log(StatSnap(reporter, true))
}
