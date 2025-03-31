package web

import (
	"math"
	"time"
)

type TaskListResponse struct {
	Tasks []Task0 `json:"tasks"` // API 列表
	Total int     `json:"total"` // API 总数
}

type TaskReq struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	AIds       []int64 `json:"a_ids"`       // 接口里可能包含 http, 也可能是 websocket
	Durations  string  `json:"durations"`   // 持续时间
	Workers    int64   `json:"workers"`     // 并发数
	MaxWorkers int64   `json:"max_workers"` // 最大持续时间
	Rate       float64 `json:"rate"`        // 速率
	Execute    bool    `json:"execute"`     // 运行性能测试
}

const (
	DefaultWorkers    uint64  = 5
	DefaultMaxWorkers uint64  = math.MaxUint64
	DefaultDurations          = 10 * time.Minute
	DefaultRate       float64 = 10
)

type Task0 struct {
	Id         int64    `json:"id"`
	Name       string   `json:"name"`
	AIds       []string `json:"a_ids"` // 接口里可能包含 http, 也可能是 websocket
	APIs       []string `json:"apis"`
	Durations  string   `json:"durations"`   // 持续时间/超时时间
	Workers    uint64   `json:"workers"`     // 并发数
	MaxWorkers uint64   `json:"max_workers"` // 最大持续时间
	Rate       float64  `json:"rate"`        // Rate 速率

	Creator string `json:"creator"`
	Updater string `json:"updater"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}

type TaskAPI0 struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
