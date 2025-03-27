package execution

//
//
//import (
//	model "TestCopilot/TestEngine/internal/service/core/model"
//	"TestCopilot/TestEngine/pkg/logger"
//	"context"
//	"fmt"
//	rate2 "golang.org/x/time/rate"
//	"runtime"
//	"sync"
//	"time"
//)
//
//type ExecutionLoadRun interface {
//	HttpRunDebug(result chan []*model.HttpResult, wg *sync.WaitGroup, s *model.Subtask) []*model.HttpResult
//	HttpRun(duration time.Duration, rate float64)
//}
//
//type executionLoadRun struct {
//	t       *model.TaskService
//	subtask model.Subtask
//	l       logger.LoggerV1
//	//report model.ReportService
//}
//
//func NewExecutionLoadRun(t *model.TaskService) ExecutionLoadRun {
//	return &executionLoadRun{
//		t: t,
//		//l: l,
//		//report: report,
//	}
//}
//
//func (e *executionLoadRun) HttpRunDebug(result chan []*model.HttpResult, wg *sync.WaitGroup, s *model.Subtask) []*model.HttpResult {
//	defer wg.Done()
//	res := make([]*model.HttpResult, 0)
//	for _, api := range e.t.APIs {
//		api_res := api.Http.Send(s)
//		api_res.Task = e.t.Name
//		res = append(res, api_res)
//	}
//	// 一次任务的结果
//	model.TaskDebugLogs(true, res)
//
//	result <- res
//	return res
//}
//
//func (e *executionLoadRun) HttpRun(duration time.Duration, rate float64) {
//	// 这里将 task 中的所有接口，按照一个goroutine 去请求，
//	// 创建限速器
//	limiter := rate2.NewLimiter(rate2.Limit(rate), 1)
//
//	// 创建上下文，用于控制压测持续时间
//	ctx, cancel := context.WithTimeout(context.Background(), duration)
//	defer cancel()
//
//	var wg sync.WaitGroup
//
//	worker := e.t.TaskConf.Workers
//	if worker > e.t.TaskConf.MaxWorkers {
//		worker = e.t.TaskConf.MaxWorkers
//	}
//
//	results := make(chan []*model.HttpResult)
//
//	s := &model.Subtask{
//		Began: time.Now(),
//	}
//
//	for i := uint64(0); i < worker; i++ {
//		wg.Add(1)
//		go e.HttpRunDebug(results, &wg, s)
//	}
//
//	go func() {
//		defer func() {
//			wg.Wait()
//			close(results)
//		}()
//		// 启动多个 goroutine 发送请求
//		for {
//			select {
//			case <-ctx.Done():
//				// 压测时间到，退出循环
//				fmt.Println("压测结束")
//				return
//			default:
//				// 判断是否获取到令牌
//				if limiter.Allow() {
//					// 启动 goroutine 发送请求
//					wg.Add(1)
//					go e.HttpRunDebug(results, &wg, s)
//				}
//				//} else {
//				//	// 处理限速，例如记录日志或等待一段时间
//				//	e.l.Warn(fmt.Sprintf("限速，当前的 rate limit 为：%v\n", limiter.Limit()))
//				//	time.Sleep(10 * time.Millisecond)
//			}
//		}
//	}()
//
//	model.FinalReport(s, results)
//	e.l.Info(fmt.Sprintf("并发请求数：%d\n", worker))
//	e.l.Info(fmt.Sprintf("当前 http_load 的 goroutine 数量: %d\n", runtime.NumGoroutine()))
//
//}
//
////func displayTask(e *executionLoadRun) string {
////	return fmt.Sprintf(`
////+++++ taskService Struct HttpResult: +++++
////[Name: %s]
////[APIs_profile: %v]
////[Durations: %v]
////[Workers: %v]
////[MaxWorkers:%v]
////[Timeout: %v]
////`, e.t.Name, e.t.APIs, e.t.TaskConf.Durations, e.t.TaskConf.Workers, e.t.TaskConf.MaxWorkers, e.t.TaskConf.Timeout)
////}
