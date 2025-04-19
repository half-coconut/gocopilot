package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	rate2 "golang.org/x/time/rate"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type TaskService interface {
	InterfacesDebug(ctx context.Context, tid int64) []TaskDebugLog
	PerformanceRun(ctx context.Context, tid int64) string
	PerformanceDebug(ctx context.Context, tid int64, result chan []*HttpResult, wg *sync.WaitGroup) []*HttpResult

	DebugForAPI(ctx context.Context, task domain.Task) TaskDebugLog

	Save(ctx *gin.Context, task domain.Task, uid int64) (int64, error)
	List(ctx context.Context, uid int64) ([]domain.Task, error)
	GetDetailByTid(ctx context.Context, tid int64) (domain.Task, error)
	SetBegin(ctx context.Context)
}

// taskService 任务结构体
type taskService struct {
	repo    repository.TaskRepository
	httpSvc HttpService
	l       logger.LoggerV1
	subtask *Subtask
}

func (t *taskService) GetDetailByTid(ctx context.Context, tid int64) (domain.Task, error) {
	return t.repo.FindByTId(ctx, tid)
}

type Subtask struct {
	Began time.Time
	seqmu sync.Mutex
	seq   uint64
}

func NewTaskService(repo repository.TaskRepository, l logger.LoggerV1, httpSvc HttpService) TaskService {
	return &taskService{
		repo:    repo,
		l:       l,
		subtask: &Subtask{},
		httpSvc: httpSvc,
	}
}

func (t *taskService) List(ctx context.Context, uid int64) ([]domain.Task, error) {
	return t.repo.FindByUId(ctx, uid)
}

func (t *taskService) SetBegin(ctx context.Context) {
	t.subtask.Began = time.Now()
}

func (t *taskService) Save(ctx *gin.Context, task domain.Task, uid int64) (int64, error) {
	if task.Id > 0 {
		// 这里是修改
		task.Updater = domain.Editor{
			Id: uid,
		}
		err := t.repo.Update(ctx, task)
		if err != nil {
			t.l.Warn("修改失败", logger.Error(err))
		}
		return task.Id, err
	}
	// 这里是新增
	task.Creator = domain.Editor{
		Id: uid,
	}
	task.Updater = domain.Editor{
		Id: uid,
	}
	Id, err := t.repo.Create(ctx, task)
	if err != nil {
		t.l.Warn("新增失败", logger.Error(err))
	}
	return Id, err
}

func (t *taskService) InterfacesDebug(ctx context.Context, tid int64) []TaskDebugLog {
	t.SetBegin(ctx)

	task := t.getTask(ctx, tid)
	var (
		reports []TaskDebugLog
		res     *HttpResult
	)
	for _, api := range task.APIs {
		if api.Type == "http" {
			headers := http.Header{}
			for key, value := range api.Header {
				headers.Add(key, value)
			}

			t.httpSvc.SetHttpInput(api.Method,
				api.URL, api.Params,
				[]byte(jsonx.JsonMarshal(api.Body)),
				headers)

			res = t.httpSvc.Send(t.subtask)

			res.Task = task.Name
			content := ExportDebugLogs(true, res)
			reports = append(reports, content)
		}
	}
	// 一次任务的结果
	return reports
}

// DebugForAPI 这个接口是为了web api 服务，后期可以使用 InterfaceDebug
func (t *taskService) DebugForAPI(ctx context.Context, task domain.Task) TaskDebugLog {
	t.SetBegin(ctx)

	api := task.APIs[0]
	var res *HttpResult
	if api.Type == "http" {
		headers := http.Header{}
		for key, value := range api.Header {
			headers.Add(key, value)
		}

		t.httpSvc.SetHttpInput(api.Method,
			api.URL, api.Params,
			[]byte(jsonx.JsonMarshal(api.Body)),
			headers,
		)

		res = t.httpSvc.Send(t.subtask)

		res.Task = task.Name
	}

	// 一次任务的结果
	content := ExportDebugLogs(true, res)
	return content
}

func (t *taskService) PerformanceDebug(ctx context.Context, tid int64, result chan []*HttpResult, wg *sync.WaitGroup) []*HttpResult {
	defer wg.Done()

	task := t.getTask(ctx, tid)

	res := make([]*HttpResult, 0)
	for _, api := range task.APIs {
		if api.Type == "http" {
			headers := http.Header{}
			for key, value := range api.Header {
				headers.Add(key, value)
			}

			t.httpSvc.SetHttpInput(api.Method, api.URL, api.Params, []byte(jsonx.JsonMarshal(api.Body)), headers)
			apiRes := t.httpSvc.Send(t.subtask)

			apiRes.Task = task.Name
			res = append(res, apiRes)
		}
	}
	// 一次任务的结果
	TaskDebugLogs(true, res)

	result <- res
	return res
}

func (t *taskService) PerformanceRun(ctx context.Context, tid int64) string {
	t.SetBegin(ctx)
	// 这里将 task 中的所有接口，按照一个goroutine 去请求，
	// 创建限速器

	task := t.getTask(ctx, tid)
	limiter := rate2.NewLimiter(rate2.Limit(task.Rate), 1)

	// 创建上下文，用于控制压测持续时间
	ctx, cancel := context.WithTimeout(context.Background(), task.Durations)
	defer cancel()

	var wg sync.WaitGroup

	worker := task.Workers
	if worker > task.MaxWorkers {
		worker = task.MaxWorkers
	}

	results := make(chan []*HttpResult)

	for i := uint64(0); i < worker; i++ {
		if limiter.Allow() {
			wg.Add(1)
			go t.PerformanceDebug(ctx, tid, results, &wg)
		}
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
					go t.PerformanceDebug(ctx, tid, results, &wg)
				}
				//} else {
				//	// 处理限速，例如记录日志或等待一段时间
				//	e.l.Warn(fmt.Sprintf("限速，当前的 rate limit 为：%v\n", limiter.Limit()))
				//	time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	content := FinalReport(t.subtask.Began, results)
	t.l.Info(fmt.Sprintf("并发请求数：%d\n", worker))
	t.l.Info(fmt.Sprintf("当前 http_load 的 goroutine 数量: %d\n", runtime.NumGoroutine()))
	return content

}

func (t *taskService) getTask(ctx context.Context, tid int64) domain.Task {
	task, err := t.repo.FindByTId(ctx, tid)
	if err != nil {
		t.l.Warn("查询任务失败", logger.Error(err))
	}
	return task
}

func TaskOnceDebugLogs(debug bool, re *HttpResult) string {
	if debug {
		content := fmt.Sprintf(`
+++++ taskService InterfacesDebug Log: +++++
[taskService: %s]
[Code: %d]
[Method: %s]
[URL:%s]
[Duration: %v]
[Headers: %v]
[Request: %s]
[Response: %s]
[Client IP: %s]
[Error: %s]`, re.Task, re.Code, re.Method, re.URL, re.Duration, re.Headers, re.Req, re.Resp, re.ClientIp, re.Error)
		log.Println(content)
		return content
	}
	return ""
}

func ExportDebugLogs(debug bool, re *HttpResult) TaskDebugLog {
	if debug {
		return TaskDebugLog{
			Name:     re.Task,
			Code:     re.Code,
			Method:   re.Method,
			Url:      re.URL,
			Duration: round(re.Duration),
			Headers:  jsonx.JsonMarshal(re.Headers),
			Request:  re.Req,
			Response: re.Resp,
			ClientIP: re.ClientIp,
			Error:    re.Error,
		}
	}
	return TaskDebugLog{}
}

type TaskDebugLog struct {
	Name     string        `json:"name"`
	Code     int64         `json:"code"`
	Method   string        `json:"method"`
	Url      string        `json:"url"`
	Duration time.Duration `json:"duration"`
	Headers  string        `json:"headers"`
	Request  string        `json:"request"`
	Response string        `json:"response"`
	ClientIP string        `json:"client_ip"`
	Error    string        `json:"error"`
}
