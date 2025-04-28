package web

type CronJobReq struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`                  // 任务名称
	Description string `json:"description,omitempty"` // 任务描述
	Type        string `json:"type"`                  // 类型: 定时任务(短任务)，持续任务(长任务)
	Cron        string `json:"cron"`                  // 定时任务: Cron "*/1 * * * *" 表达式
	HttpCfg     string `json:"http_cfg"`              // HTTP请求：endpoint 和 method
	TaskId      int64  `json:"task_id"`               // 测试任务 ID，按照 svc 的内部方法直接调用，注意一次任务执行时间和定时任务的时间不要冲突
	TimeZone    string `json:"timezone,omitempty"`    // 时区?
	Duration    string `json:"duration"`              // 持续任务: 运行时间，超时退出，用于 http 请求
	Retry       bool   `json:"retry,omitempty"`       // 是否重试
	MaxRetries  uint64 `json:"max_retries"`           // 最大重试次数
}

type CronJob0 struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`                  // 任务名称
	Description string `json:"description,omitempty"` // 任务描述
	Type        string `json:"type"`                  // 类型: 定时任务(短任务)，持续任务(长任务)
	Cron        string `json:"cron"`                  // 定时任务: Cron "*/1 * * * *" 表达式
	HttpCfg     string `json:"http_cfg"`              // HTTP请求：endpoint 和 method
	TaskId      int64  `json:"task_id"`               // 测试任务 ID，按照 svc 的内部方法直接调用
	TimeZone    string `json:"timezone,omitempty"`    // 时区?
	Duration    string `json:"duration"`              // 持续任务: 运行时间，可能用于 http 请求
	Retry       bool   `json:"retry,omitempty"`       // 是否重试
	MaxRetries  uint64 `json:"max_retries"`           // 最大重试次数
	NextTime    string `json:"next_time"`
	Status      int    `json:"status"`

	Creator string `json:"creator"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}

type CronJobListResponse struct {
	Cronjob []CronJob0 `json:"cronjob"` // Cronjob 列表
	Total   int        `json:"total"`   // Cronjob 总数
}
