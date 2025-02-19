package main

import (
	"egg_yolk/cat/log"
	"fmt"
	"sort"
	"time"
)

func (r *Report) Report(b *Base) {
	r.Requests(b)
	r.Latencies(b)
	log.L.Info(r.displayReport())
}

func (r *Report) Requests(b *Base) {
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

func (r *Report) Latencies(b *Base) {
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

func index(percent float64, b []time.Duration) time.Duration {
	idx := int64(percent / 100.0 * float64(len(b)))
	// 防止越界
	if idx > int64(len(b)) {
		idx = int64(len(b) - 1)
	}
	return b[idx]
}

var durations = [...]time.Duration{
	time.Hour,
	time.Minute,
	time.Second,
	time.Millisecond,
	time.Microsecond,
	time.Nanosecond,
}

func round(d time.Duration) time.Duration {
	for i, unit := range durations {
		if d >= unit && i < len(durations)-1 {
			return d.Round(durations[i+1])
		}
	}
	return d
}

func (r *Report) displayReport() string {
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

func printCh(resCh chan []*Result) {
	go func() {
		for {
			res, ok := <-resCh
			for _, r := range res {
				fmt.Println("+++resCh+++: ", r)
			}
			if !ok {
				break
			}
		}
	}()
}

// DisplayReport 生成 Report Base，并输出 Report
func DisplayReport(s *subtask, resCh chan []*Result) {
	var b Base
	b.Codes = make([]int, 0)
	b.SuccessRequests = make([]int, 0)
	b.FailedRequests = make([]int, 0)
	b.Durations = make([]time.Duration, 0)

	for res := range resCh {
		for _, r := range res {
			// 暂定200为成功状态码
			if r.Code == 200 {
				b.SuccessRequests = append(b.SuccessRequests, r.Code)
			} else {
				// 还没有对错误码作分类
				b.FailedRequests = append(b.FailedRequests, r.Code)
			}
			b.Codes = append(b.Codes, r.Code)
			b.Durations = append(b.Durations, r.Duration)
		}
	}

	b.TotalDuration = time.Since(s.began)

	b.TotalRequest = len(b.Codes)

	displayReportBase(b)
	var r Report
	r.Report(&b)

}

func DisplayReportResult(s *subtask, res []*Result) {
	var b Base
	b.Codes = make([]int, 0)
	b.SuccessRequests = make([]int, 0)
	b.FailedRequests = make([]int, 0)
	b.Durations = make([]time.Duration, 0)

	for _, r := range res {
		// 暂定200为成功状态码
		if r.Code == 200 {
			b.SuccessRequests = append(b.SuccessRequests, r.Code)
		} else {
			// 还没有对错误码作分类
			b.FailedRequests = append(b.FailedRequests, r.Code)
		}
		b.Codes = append(b.Codes, r.Code)
		b.Durations = append(b.Durations, r.Duration)
	}

	b.TotalDuration = time.Since(s.began)

	b.TotalRequest = len(b.Codes)

	var r Report
	r.Report(&b)
}

func displayReportBase(b Base) {
	log.L.Info(fmt.Sprintf(`
+++++ Report Base: +++++
[Codes: %v]
[TotalRequest: %d]
[SuccessRequests: %v]
[FailedRequests: %v]
[TotalDuration: %s]
[Durations:%v]
`, b.Codes, b.TotalRequest, b.SuccessRequests, b.FailedRequests, b.TotalDuration, b.Durations))
}

// %.2f: 保留两位小数，例如输出 3.14。
// %e: 使用科学计数法表示，例如输出 3.141590e+00。
// %g: 根据数值的大小自动选择 %e 或 %f 格式。

// %s 格式化字符串会将 time.Duration 类型的值转换为字符串，例如 10s、 1m30s、 2h5m10s 等。
// time.Duration 类型的值会自动转换为合适的单位，例如秒、分钟、小时等。
// %v: 等同于 %s， 输出可读的字符串格式。
