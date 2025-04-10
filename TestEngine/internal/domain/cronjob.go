package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type CronJob struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`                  // 任务名称
	Description string        `json:"description,omitempty"` // 任务描述
	Type        string        `json:"type"`                  // 类型: 定时任务(短任务)，持续任务(长任务)
	Cron        string        `json:"cron"`                  // 定时任务: Cron "*/1 * * * *" 表达式
	HttpCfg     string        `json:"http_cfg"`              // HTTP请求：endpoint 和 method
	TaskId      int64         `json:"task_id"`               // 测试任务 ID，按照 svc 的内部方法直接调用
	TimeZone    string        `json:"timezone,omitempty"`    // 时区?
	Duration    time.Duration `json:"duration"`              // 持续任务: 运行时间，可能用于 http 请求
	Retry       bool          `json:"retry,omitempty"`       // 是否重试
	MaxRetries  uint64        `json:"maxRetries"`            // 最大重试次数
	NextTime    time.Time     `json:"next_time"`

	Creator Editor `json:"creator"`
	Ctime   time.Time
	Utime   time.Time
	Version int
}

var cronjob_parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func (j CronJob) SetNextTime() time.Time {
	s, _ := cronjob_parser.Parse(j.Cron)
	return s.Next(time.Now())
}

func NextTimeV1(cron string) time.Time {
	s, _ := cronjob_parser.Parse(cron)
	return s.Next(time.Now())
}
