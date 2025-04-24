package core

import (
	"context"
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	events "github.com/half-coconut/gocopilot/core-engine/internal/events/report"
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
	GenerateSummary(ctx context.Context, begin time.Time, resCh chan []*domain.HttpResult, debug bool) string
}

type reportServiceImpl struct {
	l        logger.LoggerV1
	repo     repository.ReportRepository
	producer events.DebugLogProducer
}

func NewReportService(l logger.LoggerV1, repo repository.ReportRepository, producer events.DebugLogProducer) ReportService {
	return &reportServiceImpl{
		l:        l,
		repo:     repo,
		producer: producer,
	}
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
		go func() {
			er := svc.producer.ProducerRecordDebugLogsEvent(ctx, logs)
			if er != nil {
				svc.l.Error("发送记录debug日志事件失败")
			}
		}()

		//rid, err := svc.repo.CreateDebugLog(ctx, logs)
		//if err != nil {
		//	svc.l.Info(fmt.Sprintf("保存 Debug日志失败，rid: %v", rid), logger.Error(err))
		//	return domain.DebugLog{}, err
		//}
		return logs, nil
	}
	return domain.DebugLog{}, nil
}

func (svc *reportServiceImpl) CreateDebugLogs(ctx context.Context, debug bool, res []*domain.HttpResult) ([]domain.DebugLog, error) {
	// 这里的res 是一次任务里包含的所有接口，如果 10-20,这里就是 10-20个的 debug 信息
	var err error
	batchRes := make([]domain.DebugLog, 0)
	if debug {
		// 这里可以后期扩展为批量发送
		for _, re := range res {
			content, err := svc.CreateDebugLog(ctx, debug, re)
			err = err
			log.Println(content)
			batchRes = append(batchRes, content)
		}
	}
	return batchRes, err
}

func (svc *reportServiceImpl) GenerateSummary(ctx context.Context, begin time.Time, resCh chan []*domain.HttpResult, debug bool) string {
	// 生成 Report String
	b := svc.generateBase(begin, resCh)
	r := svc.requests(b)
	svc.l.Info(fmt.Sprintf("svc 里打印 debug 是什么：%v", debug))
	r.Debug = debug
	r = svc.latencies(r, b)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	summary, err := svc.repo.CreateSummary(ctxTimeout, *r)
	if err != nil {
		svc.l.Info(fmt.Sprintf("创建 summary 失败：%v，err: %v", summary, err))
	}
	cancel()
	svc.l.Info(fmt.Sprintf("summary 结构体：%v", summary))
	return svc.displayReport(r)
}

func (svc *reportServiceImpl) generateBase(begin time.Time, resCh chan []*domain.HttpResult) *domain.Base {
	var b domain.Base
	b.Codes = make([]int64, 0)
	b.SuccessRequests = make([]int64, 0)
	b.FailedRequests = make([]int64, 0)
	b.Durations = make([]time.Duration, 0)

	for res := range resCh {
		for _, re := range res {
			b.TaskId = re.TaskId
			b.TName = re.TName
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
	return &b
}

func (svc *reportServiceImpl) requests(b *domain.Base) *domain.Summary {
	var r domain.Summary
	r.TaskId = b.TaskId
	r.TName = b.TName
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
	return &r
}

func (svc *reportServiceImpl) latencies(r *domain.Summary, b *domain.Base) *domain.Summary {
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
	return r
}

func (svc *reportServiceImpl) displayReport(r *domain.Summary) string {
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
