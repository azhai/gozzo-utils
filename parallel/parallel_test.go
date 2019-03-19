package parallel

import (
	"fmt"
	"testing"
	"time"
)

func SetUpAction(id int) {
	time.Sleep(time.Duration(id*100) * time.Millisecond)
	fmt.Println("SetUp", id)
}

func LoopAction(id int, n uint64) {
	fmt.Println("Loop", id, n)
}

func CreateWork(qChan chan bool, msecs int) *AlignWork {
	task1 := NewStageTask(SetUpAction, false)
	task2 := NewRangeTask(LoopAction, qChan, 1)
	task2.EachGap = msecs
	task2.Config = func(id int) (uint64, uint64) {
		var stop = uint64(id)*100 + 1
		return stop - 100, stop
	}
	return new(AlignWork).Then(task1).Then(task2)
}

func RunScheduler(count, runSecs, gapMsecs int) {
	elapse := time.Duration(runSecs) * time.Second
	sch := NewScheduler(elapse, true)
	work := CreateWork(sch.QuitChan, gapMsecs)
	sch.Run(work, count)
}

func TestRun(t *testing.T) {
	RunScheduler(10, 15, 50)
}

func BenchmarRun(b *testing.B) {
	RunScheduler(b.N, 0, 1)
}
