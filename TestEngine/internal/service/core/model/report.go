package model

import (
	"TestCopilot/TestEngine/pkg/logger"
	"fmt"
	"log"
	"sort"
	"time"
)

// ReportService debug 报告，最终报告
type ReportService interface {
}

type reportService struct {
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
	l             logger.LoggerV1
}

type Base struct {
	Codes           []int
	TotalRequest    int
	SuccessRequests []int
	FailedRequests  []int
	TotalDuration   time.Duration
	Durations       []time.Duration
}

//func NewReport(l logger.LoggerV1) ReportService {
//	return &reportService{
//		l: l,
//	}
//}

// FinalReport 生成 Report Base，并输出 Report
func FinalReport(s *Subtask, resCh chan []*HttpResult) string {
	var b Base
	b.Codes = make([]int, 0)
	b.SuccessRequests = make([]int, 0)
	b.FailedRequests = make([]int, 0)
	b.Durations = make([]time.Duration, 0)

	for res := range resCh {
		for _, re := range res {
			// 暂定200为成功状态码
			if re.Code == 200 {
				b.SuccessRequests = append(b.SuccessRequests, re.Code)
			} else {
				// 还没有对错误码作分类
				b.FailedRequests = append(b.FailedRequests, re.Code)
			}
			b.Codes = append(b.Codes, re.Code)
			b.Durations = append(b.Durations, re.Duration)
		}
	}

	b.TotalDuration = time.Since(s.Began)

	b.TotalRequest = len(b.Codes)

	//r.displayReportBase(b)
	var r reportService
	return r.generateReport(&b)
}

func (r *reportService) generateReport(b *Base) string {
	r.Requests(b)
	r.Latencies(b)
	log.Println(r.displayReport())
	return r.displayReport()
}

func (r *reportService) Requests(b *Base) {
	r.Ratio = float64(len(b.SuccessRequests)) / float64(b.TotalRequest)
	r.Total = b.TotalRequest
	r.TotalDuration = b.TotalDuration

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

func (r *reportService) Latencies(b *Base) {
	sort.Slice(b.Durations, func(i, j int) bool {
		return b.Durations[i] < b.Durations[j]
	})
	r.Min = b.Durations[0]
	r.Max = b.Durations[len(b.Durations)-1]
	r.Mean = time.Duration(int64(b.TotalDuration) / int64(len(b.Durations)))

	r.P50 = index(50, b.Durations)
	r.P90 = index(90, b.Durations)
	r.P95 = index(95, b.Durations)
	r.P99 = index(99, b.Durations)
}

func (r *reportService) displayReport() string {
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
`, r.Total, r.Rate, r.Throughput,
		round(r.TotalDuration),
		round(r.Min),
		round(r.Mean),
		round(r.Max),
		round(r.P50),
		round(r.P90),
		round(r.P95),
		round(r.P99),
		r.Ratio*100, r.StatusCodes)

}

//func DisplayReportResult(s *model.Subtask, res []*model.HttpResult) {
//	var b Base
//	b.Codes = make([]int, 0)
//	b.SuccessRequests = make([]int, 0)
//	b.FailedRequests = make([]int, 0)
//	b.Durations = make([]time.Duration, 0)
//
//	for _, r := range res {
//		// 暂定200为成功状态码
//		if r.Code == 200 {
//			b.SuccessRequests = append(b.SuccessRequests, r.Code)
//		} else {
//			// 还没有对错误码作分类
//			b.FailedRequests = append(b.FailedRequests, r.Code)
//		}
//		b.Codes = append(b.Codes, r.Code)
//		b.Durations = append(b.Durations, r.Duration)
//	}
//
//	b.TotalDuration = time.Since(s.Began)
//
//	b.TotalRequest = len(b.Codes)
//
//	var r generateReport
//	r.generateReport(&b)
//}

//func (r *reportService) displayReportBase(b Base) {
//	r.l.Info(fmt.Sprintf(`
//+++++ generateReport Base: +++++
//[Codes: %v]
//[TotalRequest: %d]
//[SuccessRequests: %v]
//[FailedRequests: %v]
//[TotalDuration: %s]
//[Durations:%v]
//`, b.Codes, b.TotalRequest, b.SuccessRequests, b.FailedRequests, b.TotalDuration, b.Durations))
//}

// %.2f: 保留两位小数，例如输出 3.14。
// %e: 使用科学计数法表示，例如输出 3.141590e+00。
// %g: 根据数值的大小自动选择 %e 或 %f 格式。

// %s 格式化字符串会将 time.Duration 类型的值转换为字符串，例如 10s、 1m30s、 2h5m10s 等。
// time.Duration 类型的值会自动转换为合适的单位，例如秒、分钟、小时等。
// %v: 等同于 %s， 输出可读的字符串格式。

func TaskDebugLogs(debug bool, res []*HttpResult) {
	//resList := make([]string, 0)
	if debug {
		for _, re := range res {
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
			//resList = append(resList, content)
		}
	}
	//return resList
}
