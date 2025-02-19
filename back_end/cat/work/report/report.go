package report

import (
	"egg_yolk/cat/work"
	"fmt"
	"time"
)

// 如果想要输出报告，必须获取每一次的结果，以及最终的结果

func (r *Report) Success(code []int) {

}

func (r *Report) Requests(workers int, du time.Duration) {
	r.Total = workers
	if r.Ratio == 100.0 {
		r.Rate = float64(workers) / du.Seconds()
		r.Throughput = r.Rate
	}
}

// DisplayReport TODO: 记得改回 Report
func DisplayReport(rch chan []*work.Result) {
	for res := range rch {
		for _, r := range res {
			println(r)
		}
	}
}

func displayReport(r *Report) string {
	return fmt.Sprintf(`
+++ Requests +++
[total 总请求数: %d]
[rate 请求速率: %.2f]
[throughput 吞吐量: %g]

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
[ratio 成功率: %v]
`, r.Total, r.Rate, r.Throughput, r.TotalDuration, r.Min, r.Mean, r.Max, r.P50, r.P90, r.P95, r.P99, r.Ratio)

}

// %.2f: 保留两位小数，例如输出 3.14。
// %e: 使用科学计数法表示，例如输出 3.141590e+00。
// %g: 根据数值的大小自动选择 %e 或 %f 格式。

// %s 格式化字符串会将 time.Duration 类型的值转换为字符串，例如 10s、 1m30s、 2h5m10s 等。
// time.Duration 类型的值会自动转换为合适的单位，例如秒、分钟、小时等。
// %v: 等同于 %s， 输出可读的字符串格式。
