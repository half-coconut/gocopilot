package main

import (
	"TestCopilot/TestEngine/cat/log"
	"context"
	"fmt"
	rate2 "golang.org/x/time/rate"
	"math"
	"runtime"
	"sync"
	"time"
)

// Task 是封装接口请求的结构体，Execute 是 work 模式的并发请求的接口

// Task 并发任务，单接口，多接口
type Task struct {
	Name       string        `json:"name"`
	APIs       []API         `json:"apis"`
	Durations  time.Duration `json:"durations"`
	Workers    uint64        `json:"workers"`
	MaxWorkers uint64        `json:"maxWorkers"`
	Timeout    time.Duration `json:"timeout"`
	stopch     chan struct{}
	stopOnce   sync.Once
	mu         sync.Mutex
}

const (
	DefaultWorkers    uint64 = 5
	DefaultMaxWorkers uint64 = math.MaxUint64
	DefaultDurations         = 10 * time.Minute
	DefaultTimeout           = 30 * time.Second
)

func NewTask(name string, apis []API, workers uint64) *Task {
	return &Task{
		Name:       name,
		APIs:       apis,
		Durations:  DefaultDurations,
		Workers:    workers,
		MaxWorkers: DefaultMaxWorkers,
		Timeout:    DefaultTimeout,
		stopch:     make(chan struct{}),
		stopOnce:   sync.Once{},
	}
}

// Run 使用默认模式
func (t *Task) Run(result chan []*Result, wg *sync.WaitGroup, s *subtask) []*Result {
	defer wg.Done()
	res := make([]*Result, 0)
	for _, api := range t.APIs {
		api_res := api.Send(s)
		api_res.Task = t.Name
		res = append(res, api_res)
	}
	// 一次任务的结果
	displayTaskDebugLogs(true, res)

	result <- res
	return res
}

func (t *Task) RunV2(result chan []*Result, wg *sync.WaitGroup, s *subtask, ctx context.Context) chan []*Result {
	defer wg.Done()

	res := make([]*Result, 0)
	for _, api := range t.APIs {
		api_res := api.Send(s)
		api_res.Task = t.Name
		res = append(res, api_res)
	}

	// 一次任务的结果
	displayTaskDebugLogs(true, res)

	select {
	case <-ctx.Done():
		// 压测结束，不再发送请求
		return result
	case result <- res:
		return result
	}
}

func (t *Task) DefaultRun(maxGoroutines int, apis []API) {
	var wg sync.WaitGroup

	worker := t.Workers
	if worker > t.MaxWorkers {
		worker = t.MaxWorkers
	}

	results := make(chan []*Result)

	s := &subtask{
		began: time.Now(),
	}

	for i := uint64(0); i < worker; i++ {
		wg.Add(1)
		go t.Run(results, &wg, s)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	DisplayReport(s, results)
	log.L.Info(fmt.Sprintf("并发请求数：%d\n", worker))
}

type subtask struct {
	began time.Time
	seqmu sync.Mutex
	seq   uint64
}

func (t *Task) http_load(duration time.Duration, rate float64) {
	// 这里将 task 中的所有接口，按照一个goroutine 去请求，
	// 假设 rate=50，len(apis)=2，则实际的 rate 小于100.
	// 考虑拆分 api，将设定的 rate 等同于实际的 rate.

	// 创建限速器
	limiter := rate2.NewLimiter(rate2.Limit(rate), 1)

	// 创建上下文，用于控制压测持续时间
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup

	worker := t.Workers
	if worker > t.MaxWorkers {
		worker = t.MaxWorkers
	}

	results := make(chan []*Result)

	s := &subtask{
		began: time.Now(),
	}

	for i := uint64(0); i < worker; i++ {
		wg.Add(1)
		go t.Run(results, &wg, s)
	}

	go func() {
		defer func() {
			wg.Wait()
			close(results)
		}()
		// 启动多个 goroutine 发送请求
		for {
			select {
			case <-ctx.Done():
				// 压测时间到，退出循环
				fmt.Println("压测结束")
				return
			default:
				// 判断是否获取到令牌
				if limiter.Allow() {
					// 启动 goroutine 发送请求
					wg.Add(1)
					go t.RunV2(results, &wg, s, ctx)
				}
				//} else {
				//	// 处理限速，例如记录日志或等待一段时间
				//	log.L.Warn(fmt.Sprintf("限速，当前的 rate limit 为：%v\n", limiter.Limit()))
				//	time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	DisplayReport(s, results)
	log.L.Info(fmt.Sprintf("并发请求数：%d\n", worker))
	log.L.Info(fmt.Sprintf("当前 http_load 的 goroutine 数量: %d\n", runtime.NumGoroutine()))

}

func displayTaskDebugLogs(debug bool, r []*Result) {
	if debug {
		for _, res := range r {
			log.L.Info(fmt.Sprintf(`
+++++ Task Debug Log: +++++
[Task: %s]
[Code: %d]
[Method: %s]
[URL:%s]
[Duration: %v]
[Headers: %v]
[Body: %s]
[Client IP: %s]
[Error: %s]`, res.Task, res.Code, res.Method, res.URL, res.Duration, res.Headers, res.Body, res.ClientIp, res.Error))
		}
	}
}

func displayTask(t *Task) string {
	return fmt.Sprintf(`
+++++ Task Struct Result: +++++
[Name: %s]
[APIs_profile: %v]
[Durations: %v]
[Workers: %v]
[MaxWorkers:%v]
[Timeout: %v]
`, t.Name, t.APIs, t.Durations, t.Workers, t.MaxWorkers, t.Timeout)
}

func (t *Task) Stop() bool {
	select {
	case <-t.stopch:
		return false
	default:
		t.stopOnce.Do(func() {
			close(t.stopch)
		})
		return true
	}
}
