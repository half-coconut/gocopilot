package hit

import (
	"net/http"
)

// Result 结果
type Result struct {
	Code     int         `json:"code,string"`
	Error    string      `json:"error"`
	Body     string      `json:"body"`
	Method   string      `json:"method"`
	URL      string      `json:"url"`
	Headers  http.Header `json:"headers"`
	Duration string      `json:"duration"`
}

// 报告 总体请求数，响应时间，平均响应时间，
// 1 根据 duration 完成并发数和请求发起，结果返回
