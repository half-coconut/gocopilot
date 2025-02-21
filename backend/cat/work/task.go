package main

import (
	"context"
	"egg_yolk/cat/log"
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
	DefaultWorkers uint64 = 5
	//DefaultMaxWorkers uint64 = 10
	DefaultMaxWorkers uint64 = math.MaxUint64
	DefaultDurations         = 10 * time.Minute
	DefaultTimeout           = 30 * time.Second
)

func NewTask(name string, apis []API, workers uint64) *Task {
	return &Task{
		Name:      name,
		APIs:      apis,
		Durations: DefaultDurations,
		//Workers:    DefaultWorkers,
		Workers:    workers,
		MaxWorkers: DefaultMaxWorkers,
		Timeout:    DefaultTimeout,
		stopch:     make(chan struct{}),
		stopOnce:   sync.Once{},
	}
}

func (t *Task) Execute(result chan []*Result, s *subtask) {
	// 多个接口，按照一个任务，执行发送
	res := make([]*Result, 0)

	for _, api := range t.APIs {
		api_res := api.Send(s)
		api_res.Task = t.Name
		res = append(res, api_res)
	}

	result <- res

	// 一次任务的结果
	displayTaskDebugLogs(true, res)
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

func (t *Task) RunV3(result chan []*Result, wg *sync.WaitGroup, s *subtask) {
	defer wg.Done()

	// 发送 HTTP 请求
	res := make([]*Result, 0)
	for _, api := range t.APIs {
		api_res := api.Send(s)
		api_res.Task = t.Name
		res = append(res, api_res)
	}
	// 一次任务的结果
	displayTaskDebugLogs(true, res)
	if res != nil {
		result <- res
	}

	//if res != nil {
	//	select {
	//	case t.ticks <- struct{}{}:
	//		count++
	//	case <-t.stopch:
	//		return
	//	default:
	//		result <- res
	//	}
	//}

}

// WorkRun 使用无缓冲通道
// WorkRun 后期按照 Duration 执行时间，需要支持动态增加 wg 的数量，直到运行到 Duration 结束
// result 通道返回阻塞，待解决
func WorkRun(maxGoroutines int, name string, apis []API) {

	p := New(maxGoroutines)
	var wg sync.WaitGroup
	workers := maxGoroutines * len(apis)

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		task_worker := &Task{
			Name:       name,
			APIs:       apis,
			Durations:  DefaultDurations,
			Workers:    DefaultWorkers,
			MaxWorkers: DefaultMaxWorkers,
			Timeout:    DefaultTimeout,
		}
		go func() {
			p.Run(task_worker)
			wg.Done()
			fmt.Println("当前work内部 goroutine 数量:", runtime.NumGoroutine())
		}()

	}

	go func() {
		wg.Wait()
		p.Shutdown()
	}()

	//res, _ := <-result
	//DisplayReportResult(res)

	//FinalReport(s,result)
	log.L.Info(fmt.Sprintf("Total 总请求数：%d\n", workers))
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

type subtask struct {
	began time.Time
	seqmu sync.Mutex
	seq   uint64
}

func (t *Task) ConstantDefaultRun(maxGoroutines int, apis []API, du time.Duration) chan []*Result {
	var wg sync.WaitGroup
	worker := uint64(maxGoroutines * len(apis))

	s := &subtask{
		began: time.Now(),
	}

	results := make(chan []*Result)
	ticks := make(chan struct{})

	for i := uint64(0); i < worker; i++ {
		wg.Add(1)
		go t.Run(results, &wg, s)
	}

	go func() {
		defer func() {
			wg.Wait()
			close(results)
			t.Stop()
		}()

		count := uint64(0)
		for {
			elapsed := time.Since(s.began)
			if du > 0 && elapsed > du {
				return
			}

			if worker < t.MaxWorkers {
				select {
				case ticks <- struct{}{}:
					count++
					continue
				case <-t.stopch:
					return
				default:
					// all workers are blocked. start one more and try again
					worker++
					wg.Add(1)
					go t.Run(results, &wg, s)
				}
			}

			select {
			case ticks <- struct{}{}:
				count++
			case <-t.stopch:
				return
			}
		}
	}()

	//go func() {
	//wg.Wait()
	//close(results)
	//t.Stop()
	//}()

	DisplayReport(s, results)
	log.L.Info(fmt.Sprintf("并发请求数：%d\n", worker))
	return results
}

func (t *Task) http_load(duration time.Duration, rate float64) {
	// 创建限速器
	limiter := rate2.NewLimiter(rate2.Limit(rate), 1)

	// 创建上下文，用于控制压测持续时间
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// 创建 WaitGroup，用于等待所有请求完成
	var wg sync.WaitGroup
	results := make(chan []*Result)

	s := &subtask{
		began: time.Now(),
	}

	// 启动多个 goroutine 发送请求
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			close(results)
			fmt.Println(results)
			res, _ := <-results
			fmt.Printf("+++res:%v\n", res) // []
			//for res := range results {
			//	fmt.Printf("+++results+++: %v\n", res)
			//}
			fmt.Println("当前 http_load 的 goroutine 数量:", runtime.NumGoroutine())
			// 压测时间到，退出循环
			fmt.Println("压测结束")
			return
		default:
			// 判断是否获取到令牌
			if limiter.Allow() {
				// 启动 goroutine 发送请求
				wg.Add(1)
				go func() {
					// 发送 HTTP 请求
					t.RunV2(results, &wg, s, ctx)
				}()
			}
			//} else {
			//	// 处理限速，例如记录日志或等待一段时间
			//	log.L.Warn(fmt.Sprintf("限速，当前的 rate limit 为：%v\n", limiter.Limit()))
			//	time.Sleep(10 * time.Millisecond)
		}
	}

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
[Error: %s]`, res.Task, res.Code, res.Method, res.URL, res.Duration, res.Headers, res.Body, res.Error))
		}
	}
}

func displayTask(t *Task) string {
	return fmt.Sprintf(`
+++++ Task Struct Result: +++++
[Name: %s]
[APIs: %v]
[Durations: %v]
[Workers: %v]
[MaxWorkers:%v]
[Timeout: %v]
`, t.Name, t.APIs, t.Durations, t.Workers, t.MaxWorkers, t.Timeout)
}
