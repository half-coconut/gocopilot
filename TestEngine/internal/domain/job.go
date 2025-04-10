package domain

import (
	"github.com/robfig/cron/v3"
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

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func (j Job) NextTime() time.Time {
	s, _ := parser.Parse(j.Cron)
	return s.Next(time.Now())
}
