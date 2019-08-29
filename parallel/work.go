package parallel

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

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

// Actor模式
type Actor struct {
	task    ITask
	workers int
	pool    chan int
	done    chan bool
	mu      sync.Mutex
}

func NewActor(task ITask, workers int) *Actor {
	return &Actor{
		task:    task,
		workers: workers,
		pool:    make(chan int),
		done:    make(chan bool, workers),
	}
}

func (a *Actor) Close() {
	a.mu.Lock()
	for i := 0; i < a.workers; i++ {
		a.done <- true
	}
	close(a.done)
	close(a.pool)
	a.mu.Unlock()
}

func (a *Actor) Add(id int) {
	a.mu.Lock()
	a.pool <- id
	a.mu.Unlock()
}

func (a *Actor) Work() {
	for {
		select {
		case <-a.done:
			return
		case id := <-a.pool:
			if a.task.IsRoutine() {
				go a.task.Process(id)
			} else {
				a.task.Process(id)
			}
		}
	}
}

func (a *Actor) Run() {
	for i := 0; i < a.workers; i++ {
		go a.Work()
	}
}
