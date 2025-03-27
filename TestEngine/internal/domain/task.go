package domain

import (
	"net/http"
	"time"
)

// Task 任务结构体 Domain
type Task struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	APIs []APIs `json:"apis"` // 接口里可能包含 http, 也可能是 websocket

	Durations  time.Duration `json:"durations"`  // 持续时间
	Workers    uint64        `json:"workers"`    // 并发数
	MaxWorkers uint64        `json:"maxWorkers"` // 最大持续时间
	Timeout    time.Duration `json:"timeout"`    // 超时时间

	Creator Editor `json:"creator"`
	Updater Editor `json:"updater"`
	Ctime   time.Time
	Utime   time.Time
}

type APIs struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Params string `json:"params,omitempty"`
	Body   string `json:"body,omitempty"`
	Header string `json:"header,omitempty"`
	Method string `json:"method"`
	Type   string `json:"type,omitempty"` // http/websocket
}

type TaskAPI struct {
	Name    string                 `json:"name"`
	URL     string                 `json:"url"`
	Params  string                 `json:"params"`
	Body    map[string]interface{} `json:"body"`
	Headers map[string]string      `json:"headers"`
	Method  string                 `json:"method"`
}

// HttpContent 参考
type HttpContent struct {
	Name   string      `json:"name"`
	URL    string      `json:"url"`
	Params string      `json:"params,omitempty"`
	Body   []byte      `json:"data,omitempty"`
	Header http.Header `json:"header,omitempty"`
	Method string      `json:"method"`
}
