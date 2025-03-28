package model

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository"
	"TestCopilot/TestEngine/pkg/jsonx"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	rate2 "golang.org/x/time/rate"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type TaskService interface {
	Debug(taskDomain domain.Task) TaskDebugLog
	OnceRunDebug(taskDomain domain.Task) string
	HttpRunDebug(taskDomain domain.Task, result chan []*HttpResult, wg *sync.WaitGroup, s *Subtask) []*HttpResult
	HttpRun(taskDomain domain.Task, duration time.Duration, rate float64)
	Save(ctx *gin.Context, task domain.Task, uid int64) (int64, error)
}

// taskService 任务结构体
type taskService struct {
	repo repository.TaskRepository
	l    logger.LoggerV1
}

type Subtask struct {
	Began time.Time
	seqmu sync.Mutex
	seq   uint64
}

func NewTaskService(repo repository.TaskRepository, l logger.LoggerV1) TaskService {
	return &taskService{
		repo: repo,
		l:    l,
	}
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

func (t *taskService) OnceRunDebug(taskDomain domain.Task) string {
	s := &Subtask{
		Began: time.Now(),
	}
	task := taskDomain.APIs[0]
	var res *HttpResult
	// 根据 type 区分不同的协议，需要那 api 转为 http 请求，然后发送
	if task.Type == "http" {
		var h = make(http.Header, 0)
		h.Add("Content-Type", "application/json")
		h.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")

		body := []byte(`{"jsonrpc": "2.0", "method": "eth_accounts", "params": [], "id": 1}`)
		target := NewHttpContent(
			//task.Method, task.URL, task.Params, []byte(task.Body), jsonToHeader(task.Header))
			task.Method, task.URL, task.Params, body, h)
		res = target.Send(s)
		res.Task = taskDomain.Name
	}

	// 一次任务的结果
	content := TaskOnceDebugLogs(true, res)

	return content
}

func (t *taskService) Debug(taskDomain domain.Task) TaskDebugLog {
	s := &Subtask{
		Began: time.Now(),
	}
	task := taskDomain.APIs[0]
	var res *HttpResult
	// 根据 type 区分不同的协议，需要那 api 转为 http 请求，然后发送
	if task.Type == "http" {
		headers := http.Header{}
		for key, value := range task.Header {
			headers.Add(key, value)
		}

		target := NewHttpContent(task.Method,
			task.URL, task.Params,
			[]byte(jsonx.JsonMarshal(task.Body)),
			headers,
		)
		res = target.Send(s)
		res.Task = taskDomain.Name
	}

	// 一次任务的结果
	//content := TaskOnceDebugLogs(true, res)
	content := ExportDebugLogs(true, res)
	return content
}

func (t *taskService) HttpRunDebug(taskDomain domain.Task, result chan []*HttpResult, wg *sync.WaitGroup, s *Subtask) []*HttpResult {
	defer wg.Done()

	res := make([]*HttpResult, 0)
	for _, api := range taskDomain.APIs {
		if api.Type == "http" {
			headers := http.Header{}
			for key, value := range api.Header {
				headers.Add(key, value)
			}

			target := NewHttpContent(
				api.Method, api.URL, api.Params, []byte(jsonx.JsonMarshal(api.Body)), headers)
			api_res := target.Send(s)
			api_res.Task = taskDomain.Name
			res = append(res, api_res)
		}
	}
	// 一次任务的结果
	TaskDebugLogs(true, res)

	result <- res
	return res
}

func (t *taskService) HttpRun(taskDomain domain.Task, duration time.Duration, rate float64) {
	// 这里将 task 中的所有接口，按照一个goroutine 去请求，
	// 创建限速器
	limiter := rate2.NewLimiter(rate2.Limit(rate), 1)

	// 创建上下文，用于控制压测持续时间
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup

	worker := taskDomain.Workers
	if worker > taskDomain.MaxWorkers {
		worker = taskDomain.MaxWorkers
	}

	results := make(chan []*HttpResult)

	s := &Subtask{
		Began: time.Now(),
	}

	for i := uint64(0); i < worker; i++ {
		wg.Add(1)
		go t.HttpRunDebug(taskDomain, results, &wg, s)
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
					go t.HttpRunDebug(taskDomain, results, &wg, s)
				}
				//} else {
				//	// 处理限速，例如记录日志或等待一段时间
				//	e.l.Warn(fmt.Sprintf("限速，当前的 rate limit 为：%v\n", limiter.Limit()))
				//	time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	FinalReport(s, results)
	t.l.Info(fmt.Sprintf("并发请求数：%d\n", worker))
	t.l.Info(fmt.Sprintf("当前 http_load 的 goroutine 数量: %d\n", runtime.NumGoroutine()))

}

// 从JSON字符串转换回http.Header
func jsonToHeader(headerJSON string) http.Header {
	// 创建一个用于解析的map
	var headerMap map[string][]string
	err := json.Unmarshal([]byte(headerJSON), &headerMap)
	if err != nil {
		return nil
	}
	// 将map转换为http.Header
	header := make(http.Header)
	for key, values := range headerMap {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	return header
}

func TaskOnceDebugLogs(debug bool, re *HttpResult) string {
	if debug {
		content := fmt.Sprintf(`
+++++ taskService Debug Log: +++++
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
