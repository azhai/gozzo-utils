package parallel

import (
	"time"

	"github.com/azhai/gozzo-utils/random"
)

// 间隔重复执行
func TickTime(ms int) (tick <-chan time.Time) {
	if ms > 0 {
		duration := time.Duration(ms) * time.Millisecond
		tick = time.Tick(duration)
	}
	return
}

// 延时，避免拥塞
func DelayTime(ms int) {
	if ms > 0 {
		duration := time.Duration(ms) * time.Millisecond
		time.Sleep(duration)
	}
}

// 随机延时
func DelayRand(ms int) {
	if ms > 0 {
		DelayTime(random.RandInt(ms))
	}
}

// 任务接口
type ITask interface {
	IsRoutine() bool
	Process(id int)
}

// 阶段任务
type StageTask struct {
	isRoutine bool // 是否协程中执行
	process   func(id int)
}

func NewStageTask(process func(id int), isRoutine bool) *StageTask {
	return &StageTask{process: process, isRoutine: isRoutine}
}

func (t *StageTask) IsRoutine() bool {
	return t.isRoutine
}

func (t *StageTask) Process(id int) {
	t.process(id)
}

// 区间任务
type RangeTask struct {
	QuitChan chan bool
	Action   func(id int, n uint64)
	Config   func(id int) (uint64, uint64)
	Step     uint64 // 步进，无限循环为0
	EachGap  int    // 每次的休息间隔，单位ms
	MaxDelay int    // 启动的最大延迟，单位ms
}

func NewRangeTask(action func(id int, n uint64),
	quitChan chan bool, step uint64) *RangeTask {
	return &RangeTask{Action: action, QuitChan: quitChan, Step: step}
}

func (t *RangeTask) IsRoutine() bool {
	return true
}

func (t *RangeTask) Process(id int) {
	// 计算启动和中间休眠时间
	var GapTime = time.Duration(0)
	if t.EachGap > 0 {
		GapTime = time.Duration(t.EachGap) * time.Millisecond
	}
	if t.MaxDelay > 0 {
		DelayRand(t.MaxDelay) // 随机延迟
	}
	// 计算起止范围
	var start, stop = uint64(0), uint64(1)
	if t.Config != nil {
		start, stop = t.Config(id)
	}
	// 执行循环
	for start < stop {
		select {
		default:
			t.Action(id, start)
			start += t.Step
			if GapTime > 0 {
				time.Sleep(GapTime)
			}
		case <-t.QuitChan:
			return
		}
	}
}

// 重复任务
type LoopTask struct {
	*RangeTask
}

func NewLoopTask(action func(id int, n uint64),
	quitChan chan bool) *LoopTask {
	return &LoopTask{NewRangeTask(action, quitChan, 0)}
}
