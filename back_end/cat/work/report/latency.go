package report

import (
	"fmt"
	"sort"
)

func main() {
	// 模拟一些响应时间数据 (单位：毫秒)
	responseTimes := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 200, 300, 500}

	// 计算平均响应时间
	avgResponseTime := calculateAverageResponseTime(responseTimes)
	fmt.Printf("平均响应时间: %.2f ms\n", avgResponseTime)

	// 计算不同分位数的响应时间
	p99 := calculatePercentile(responseTimes, 99)
	p95 := calculatePercentile(responseTimes, 95)
	p90 := calculatePercentile(responseTimes, 90)
	fmt.Printf("99%% 响应时间: %d ms\n", p99)
	fmt.Printf("95%% 响应时间: %d ms\n", p95)
	fmt.Printf("90%% 响应时间: %d ms\n", p90)
}

func calculateAverageResponseTime(responseTimes []int) float64 {
	total := 0
	for _, rt := range responseTimes {
		total += rt
	}
	return float64(total) / float64(len(responseTimes))
}

func calculatePercentile(responseTimes []int, percentile int) int {
	sort.Ints(responseTimes)
	index := (percentile / 100) * len(responseTimes)
	if index >= len(responseTimes) {
		index = len(responseTimes) - 1
	}
	return responseTimes[index]
}
