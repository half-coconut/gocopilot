package model

import (
	"math"
	"sync"
	"time"
)

// Task 任务结构体
type Task struct {
	Name     string     `json:"name"`
	APIs     []API      `json:"apis"`
	TaskConf TaskConfig `json:"task_config"`
	Stopch   chan struct{}
	StopOnce sync.Once
	mu       sync.Mutex
}

type TaskConfig struct {
	Durations  time.Duration `json:"durations"`  // 持续时间
	Workers    uint64        `json:"workers"`    // 并发数
	MaxWorkers uint64        `json:"maxWorkers"` // 最大持续时间
	Timeout    time.Duration `json:"timeout"`    // 超时时间
}

const (
	DefaultWorkers    uint64 = 5
	DefaultMaxWorkers uint64 = math.MaxUint64
	DefaultDurations         = 10 * time.Minute
	DefaultTimeout           = 30 * time.Second
)

func NewTaskConfig(workers uint64) *TaskConfig {
	return &TaskConfig{
		Durations:  DefaultDurations,
		Workers:    workers,
		MaxWorkers: DefaultMaxWorkers,
		Timeout:    DefaultTimeout,
	}
}

func NewTask(name string, apis []API, task_conf TaskConfig) *Task {
	return &Task{
		Name:     name,
		APIs:     apis,
		TaskConf: task_conf,
		Stopch:   make(chan struct{}),
		StopOnce: sync.Once{},
	}
}

type Subtask struct {
	Began time.Time
	seqmu sync.Mutex
	seq   uint64
}
