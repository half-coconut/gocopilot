package domain

import (
	"net/http"
	"time"
)

type DebugLog struct {
	Id       int64         `json:"id"`
	TaskId   int64         `json:"task_id"`
	AId      int64         `json:"a_id"`
	AName    string        `json:"a_name"`
	Code     int64         `json:"code"`
	Method   string        `json:"method"`
	Url      string        `json:"url"`
	Duration time.Duration `json:"duration"`
	Headers  string        `json:"headers"`
	Params   string        `json:"params"`
	Body     string        `json:"body"`
	Response string        `json:"response"`
	ClientIP string        `json:"client_ip"`
	Error    string        `json:"error"`
}

// taskId, apiId

// HttpResult 单个请求的响应结果
type HttpResult struct {
	TaskId    int64         `json:"task_id"`
	AId       int64         `json:"a_id"`
	AName     string        `json:"a_name"`
	TName     string        `json:"t_name"`
	Code      int64         `json:"code"`
	Error     string        `json:"error"`
	Params    string        `json:"params"`
	Body      string        `json:"body"`
	Resp      string        `json:"response"`
	Method    string        `json:"method"`
	URL       string        `json:"url"`
	Headers   http.Header   `json:"headers"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Seq       uint64        `json:"seq"`
	ClientIp  string        `json:"clientIp"`
}

type Summary struct {
	Id            int64
	TaskId        int64
	AIds          int64
	TName         string
	Debug         bool
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
	Status        int
}

type Base struct {
	TaskId          int64
	TName           string
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

const (
	Unknown = iota
	Passed
	Failed
	Skipped
	Errors
)
