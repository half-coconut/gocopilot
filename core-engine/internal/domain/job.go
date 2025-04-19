package domain

import (
	cronv3 "github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	// 通用的任务抽象
	Id         int64
	Name       string // 可以做唯一索引
	Cron       string
	Executor   string // http模式, LocalFuncExecutor, Schedule
	Cfg        string // http 请求任务
	Status     int
	CancelFunc func() error
}

var parser = cronv3.NewParser(cronv3.Minute | cronv3.Hour | cronv3.Dom | cronv3.Month | cronv3.Dow | cronv3.Descriptor)

func (j Job) NextTime() time.Time {
	s, _ := parser.Parse(j.Cron)
	return s.Next(time.Now())
}
