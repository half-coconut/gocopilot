package core

import (
	"context"
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"github.com/half-coconut/gocopilot/core-engine/pkg/timex"
	"log"
	"sort"
	"time"
)

// ReportService debug 报告，最终报告
type ReportService interface {
	CreateDebugLog(ctx context.Context, debug bool, re *domain.HttpResult) (domain.DebugLog, error)
	CreateDebugLogs(ctx context.Context, debug bool, res []*domain.HttpResult) ([]domain.DebugLog, error)
	GenerateReport(begin time.Time, resCh chan []*domain.HttpResult) string
}

type reportServiceImpl struct {
	l    logger.LoggerV1
	repo repository.ReportRepository
}

func NewReportService(l logger.LoggerV1, repo repository.ReportRepository) ReportService {
	return &reportServiceImpl{
		l:    l,
		repo: repo}
}

func (svc *reportServiceImpl) CreateDebugLog(ctx context.Context, debug bool, res *domain.HttpResult) (domain.DebugLog, error) {
	// 一个任务里，一个接口的 debug 信息
	if debug {
		logs := domain.DebugLog{
			TaskId:   res.TaskId,
			AId:      res.AId,
			AName:    res.AName,
			Code:     res.Code,
			Method:   res.Method,
			Url:      res.URL,
			Duration: timex.Round(res.Duration),
			Headers:  jsonx.JsonMarshal(res.Headers),
			Params:   res.Params,
			Body:     res.Body,
			Response: res.Resp,
			ClientIP: res.ClientIp,
			Error:    res.Error,
		}
		rid, err := svc.repo.CreateDebugLog(ctx, logs)
		if err != nil {
			svc.l.Info(fmt.Sprintf("保存 Debug日志失败，rid: %v", rid), logger.Error(err))
			return domain.DebugLog{}, err
		}
		return logs, nil
	}
	return domain.DebugLog{}, nil
}

func (svc *reportServiceImpl) CreateDebugLogs(ctx context.Context, debug bool, res []*domain.HttpResult) ([]domain.DebugLog, error) {
	// 这里的res 是一次任务里包含的所有接口，如果 10-20,这里就是 10-20个的 debug 信息
	// 存入数据库，还是一个任务，一个接口的存放
	var err error
	batchRes := make([]domain.DebugLog, 0)
	if debug {
		for _, re := range res {
			content, err := svc.CreateDebugLog(ctx, debug, re)
			err = err
			log.Println(content)
			batchRes = append(batchRes, content)
		}
	}
	return batchRes, err
}

func (svc *reportServiceImpl) GenerateReport(begin time.Time, resCh chan []*domain.HttpResult) string {
	// 生成 Report String
	var r Summary
	b := r.generateBase(begin, resCh)
	r.requests(&b)
	r.latencies(&b)
	return r.displayReport()
}

type Summary struct {
	Total         int
	Rate          float64
	Throughput    float64
	TotalDuration time.Duration
	Min           time.Duration
	Mean          time.Duration
	Max           time.Duration
	P50           time.Duration
	P90           time.Duration
	P95           time.Duration
	P99           time.Duration
	Ratio         float64
	StatusCodes   string
	TestStatus    TestStatus
}

type Base struct {
	Codes           []int64
	TotalRequest    int
	SuccessRequests []int64
	FailedRequests  []int64
	TotalDuration   time.Duration
	Durations       []time.Duration
	TestStatus      TestStatus
}

type TestStatus struct {
	Passed  int64
	Failed  int64
	Skipped int64
	Errors  int64
}

func (r *Summary) generateBase(begin time.Time, resCh chan []*domain.HttpResult) Base {
	var b Base
	b.Codes = make([]int64, 0)
	b.SuccessRequests = make([]int64, 0)
	b.FailedRequests = make([]int64, 0)
	b.Durations = make([]time.Duration, 0)

	for res := range resCh {
		for _, re := range res {
			// 暂定200为成功状态码
			if re.Code == int64(200) {
				b.SuccessRequests = append(b.SuccessRequests, re.Code)
				b.TestStatus.Passed += 1
			} else {
				// 还没有对错误码作分类
				b.FailedRequests = append(b.FailedRequests, re.Code)
				b.TestStatus.Failed += 1
			}
			b.Codes = append(b.Codes, re.Code)
			b.Durations = append(b.Durations, re.Duration)
		}
	}
	b.TestStatus.Skipped = 0
	b.TestStatus.Errors = 0

	b.TotalDuration = time.Since(begin)

	b.TotalRequest = len(b.Codes)
	return b
}

func (r *Summary) requests(b *Base) {
	r.Ratio = float64(len(b.SuccessRequests)) / float64(b.TotalRequest)
	r.Total = b.TotalRequest
	r.TotalDuration = b.TotalDuration
	r.TestStatus = b.TestStatus

	r.Rate = float64(r.Total) / r.TotalDuration.Seconds()
	r.Throughput = float64(len(b.SuccessRequests)) / r.TotalDuration.Seconds()

	if len(b.SuccessRequests) != 0 && len(b.FailedRequests) != 0 {
		r.StatusCodes = fmt.Sprintf(" %d:%d, %d...:%d", b.SuccessRequests[0], len(b.SuccessRequests),
			b.FailedRequests[0], len(b.FailedRequests))
	} else if len(b.SuccessRequests) != 0 && len(b.FailedRequests) == 0 {
		r.StatusCodes = fmt.Sprintf(" %d:%d", b.SuccessRequests[0], len(b.SuccessRequests))
	} else if len(b.SuccessRequests) == 0 && len(b.FailedRequests) != 0 {
		r.StatusCodes = fmt.Sprintf(" %d...:%d", b.FailedRequests[0], len(b.FailedRequests))
	}

}

func (r *Summary) latencies(b *Base) {
	sort.Slice(b.Durations, func(i, j int) bool {
		return b.Durations[i] < b.Durations[j]
	})
	r.Min = b.Durations[0]
	r.Max = b.Durations[len(b.Durations)-1]
	r.Mean = time.Duration(int64(b.TotalDuration) / int64(len(b.Durations)))

	r.P50 = timex.Index(50, b.Durations)
	r.P90 = timex.Index(90, b.Durations)
	r.P95 = timex.Index(95, b.Durations)
	r.P99 = timex.Index(99, b.Durations)
}

func (r *Summary) displayReport() string {
	return fmt.Sprintf(`
+++ Requests +++
[total 总请求数: %d]
[rate 请求速率: %.2f]
[throughput 吞吐量: %.2f]

+++ Duration +++
[total 总持续时间: %v]

+++ Latencies +++
[min 最小响应时间: %v]
[mean 平均响应时间: %v]
[max 最大响应时间: %v]
[P50 百分之50 响应时间 (中位数): %v]
[P90 百分之90 响应时间: %v]
[P95 百分之95 响应时间: %v]
[P99 百分之99 响应时间: %v]

+++ Success +++
[ratio 成功率: %.2f%%]
[status codes: %v]
[passed: %v]
[failed: %v]
`, r.Total, r.Rate, r.Throughput,
		timex.Round(r.TotalDuration),
		timex.Round(r.Min),
		timex.Round(r.Mean),
		timex.Round(r.Max),
		timex.Round(r.P50),
		timex.Round(r.P90),
		timex.Round(r.P95),
		timex.Round(r.P99),
		r.Ratio*100,
		r.StatusCodes,
		r.TestStatus.Passed,
		r.TestStatus.Failed)

}

// %.2f: 保留两位小数，例如输出 3.14。
// %e: 使用科学计数法表示，例如输出 3.141590e+00。
// %g: 根据数值的大小自动选择 %e 或 %f 格式。

// %s 格式化字符串会将 time.Duration 类型的值转换为字符串，例如 10s、 1m30s、 2h5m10s 等。
// time.Duration 类型的值会自动转换为合适的单位，例如秒、分钟、小时等。
// %v: 等同于 %s， 输出可读的字符串格式。
