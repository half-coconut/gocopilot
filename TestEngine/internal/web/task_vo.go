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
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AIds       string `json:"a_ids"`       // 接口里可能包含 http, 也可能是 websocket
	Durations  string `json:"durations"`   // 持续时间
	Workers    int64  `json:"workers"`     // 并发数
	MaxWorkers int64  `json:"max_workers"` // 最大持续时间
	Timeout    string `json:"timeout"`     // 超时时间
}

const (
	DefaultWorkers    uint64 = 5
	DefaultMaxWorkers uint64 = math.MaxUint64
	DefaultDurations         = 10 * time.Minute
	DefaultTimeout           = 30 * time.Second
)

type Task0 struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	AIds       []int64 `json:"a_ids"`       // 接口里可能包含 http, 也可能是 websocket
	Durations  string  `json:"durations"`   // 持续时间
	Workers    uint64  `json:"workers"`     // 并发数
	MaxWorkers uint64  `json:"max_workers"` // 最大持续时间
	Timeout    string  `json:"timeout"`     // 超时时间

	Creator string `json:"creator"`
	Updater string `json:"updater"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}
