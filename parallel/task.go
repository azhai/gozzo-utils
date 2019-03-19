package parallel

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
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

// 先完成全部完成Task0，再开始Task1，然后Task2
type AlignWork struct {
	tasks []ITask
}

func (w *AlignWork) Count() int {
	return len(w.tasks)
}

func (w *AlignWork) Then(task ITask) *AlignWork {
	w.tasks = append(w.tasks, task)
	return w
}

func (w *AlignWork) GetTask(index int) ITask {
	if index < 0 || index >= w.Count() {
		return nil
	}
	return w.tasks[index]
}

// 调度器
type Scheduler struct {
	waiter   *sync.WaitGroup // 同步锁
	SignChan chan os.Signal
	QuitChan chan bool
	Begin    time.Time          // 开始时间
	Elapse   time.Duration      // 超时或经历时间
	IsCtrl   bool               // 捕获Ctrl+C等系统信号
	Finally  func(s *Scheduler) //最后执行的收尾工作
}

func NewScheduler(elapse time.Duration, isCtrl bool) *Scheduler {
	runtime.GOMAXPROCS(runtime.NumCPU()) // 最多使用N个核
	if elapse <= 0 {
		elapse = 366 * 24 * time.Hour // 1年
	}
	return &Scheduler{
		waiter: new(sync.WaitGroup),
		Elapse: elapse,
		IsCtrl: isCtrl,
	}
}

func (s *Scheduler) SetFinally(finalWork func(s *Scheduler)) {
	s.Finally = finalWork
}

func (s *Scheduler) ExecTask(task ITask, count int) {
	for id := 1; id <= count; id++ {
		task.Process(id)
	}
}

func (s *Scheduler) GoExecTask(task ITask, count int) {
	s.waiter.Add(count)
	for id := 1; id <= count; id++ {
		go func(id int) {
			defer s.waiter.Done()
			task.Process(id)
		}(id)
	}
	runtime.Gosched()
	s.waiter.Wait()
}

func (s *Scheduler) Run(work *AlignWork, count int) {
	s.SetUp(count)
	for i := 0; i < work.Count(); i++ {
		task := work.GetTask(i)
		if task.IsRoutine() {
			s.GoExecTask(task, count)
		} else {
			s.ExecTask(task, count)
		}
	}
	s.TearDown()
}

func (s *Scheduler) SetUp(count int) {
	s.Begin = time.Now()
	s.QuitChan = make(chan bool, count)
	s.SignChan = make(chan os.Signal, 1)
	signal.Notify(s.SignChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-s.SignChan: // 按Ctrl+C终止
			if s.IsCtrl {
				s.TearDown()
			}
			return
		case <-time.After(s.Elapse): // 超时终止
			s.TearDown()
			return
		}
	}()
}

func (s *Scheduler) TearDown() {
	s.Elapse = time.Since(s.Begin)
	if s.QuitChan != nil {
		for i := 0; i < cap(s.QuitChan); i++ {
			s.QuitChan <- true
		}
		close(s.QuitChan)
	}
	s.Finally(s)
	os.Exit(0)
}
