package main

import (
	"net/http"
	"time"
)

// Result 单个请求的响应结果
type Result struct {
	Task      string        `json:"task"`
	Code      int           `json:"code,string"`
	Error     string        `json:"error"`
	Body      string        `json:"body"`
	Method    string        `json:"method"`
	URL       string        `json:"url"`
	Headers   http.Header   `json:"headers"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Seq       uint64        `json:"seq"`
	//Assert   Assert      `json:"assert"`
}

// Assert TODO:简单断言
type Assert struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Report struct {
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
}

type Base struct {
	Codes           []int
	TotalRequest    int
	SuccessRequests []int
	FailedRequests  []int
	TotalDuration   time.Duration
	Durations       []time.Duration
}

// Pool 和 Work 模式比较：...
