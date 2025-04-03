package core

import (
	"net/http"
	"time"
)

// HttpResult 单个请求的响应结果
type HttpResult struct {
	Task  string `json:"task"`
	Code  int64  `json:"code,string"`
	Error string `json:"error"`
	//Body      string        `json:"body"`
	Req       string        `json:"request"`
	Resp      string        `json:"response"`
	Method    string        `json:"method"`
	URL       string        `json:"url"`
	Headers   http.Header   `json:"headers"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Seq       uint64        `json:"seq"`
	ClientIp  string        `json:"clientIp"`
	//Assert   Assert      `json:"assert"`
}
