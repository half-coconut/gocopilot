package domain

import (
	"time"
)

// Task 任务结构体 Domain
type Task struct {
	Id         int64         `json:"id"`
	Name       string        `json:"name"`
	APIs       []TaskAPI     `json:"apis"` // 接口里可能包含 http, 也可能是 websocket
	AIds       []int64       `json:"a_ids"`
	Durations  time.Duration `json:"durations"`  // 持续时间
	Workers    uint64        `json:"workers"`    // 并发数
	MaxWorkers uint64        `json:"maxWorkers"` // 最大持续时间
	Rate       float64       `json:"rate"`       // rate 速率

	Creator Editor `json:"creator"`
	Updater Editor `json:"updater"`
	Ctime   time.Time
	Utime   time.Time
}

type TaskAPI struct {
	Id     int64                  `json:"id"`
	Name   string                 `json:"name"`
	URL    string                 `json:"url"`
	Params string                 `json:"params,omitempty"`
	Body   map[string]interface{} `json:"body,omitempty"`
	Header map[string]string      `json:"header,omitempty"`
	Method string                 `json:"method"`
	Type   string                 `json:"type,omitempty"` // http/websocket
}
