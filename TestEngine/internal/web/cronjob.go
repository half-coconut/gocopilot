package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/errs"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/service/core"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type CronJobHandler struct {
	l       logger.LoggerV1
	svc     service.CronJobService
	taskSvc core.TaskService
}

func NewCronJobHandler(l logger.LoggerV1, svc service.CronJobService, taskSvc core.TaskService) *CronJobHandler {
	return &CronJobHandler{l: l, svc: svc, taskSvc: taskSvc}
}

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
	MaxRetries  uint64 `json:"maxRetries"`            // 最大重试次数
}

// t.svc.PerformanceRun(ctx, req.tid)
func (c *CronJobHandler) RegisterRoutes(server *gin.Engine) {
	job := server.Group("/job")
	job.POST("/add", ginx.WrapToken[ijwt.UserClaims](c.AddAll))
	job.POST("/add/task", ginx.WrapToken[ijwt.UserClaims](c.AddInternalTask))
	job.POST("/add/http", ginx.WrapToken[ijwt.UserClaims](c.AddHttpMode))
}

// AddHttpMode 添加Http请求类型
func (c *CronJobHandler) AddHttpMode(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	var req CronJobReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.JobInvalidInput,
			Message: "用户输入格式不正确",
		}, err
	}
	duration, _ := time.ParseDuration(req.Duration)
	job := domain.CronJob{
		Id:          req.Id, // 传入了 id 就是修改，不传 id 就是新增
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Cron:        req.Cron,
		HttpCfg:     req.HttpCfg,
		TimeZone:    req.TimeZone,
		Duration:    duration,
		Retry:       req.Retry,
		MaxRetries:  req.MaxRetries,
	}

	Id, err := c.svc.Save(ctx, job, uc.Id)
	if err != nil {
		c.l.Info(fmt.Sprintf("创建 Job 失败，用户 Id：%v", uc.Id), logger.Error(err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    Id,
	}, nil
}

// AddInternalTask 添加任务类型
func (c *CronJobHandler) AddInternalTask(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	var req CronJobReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.JobInvalidInput,
			Message: "用户输入格式不正确",
		}, err
	}

	job := domain.CronJob{
		Id:          req.Id, // 传入了 id 就是修改，不传 id 就是新增
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Cron:        req.Cron,
		TaskId:      req.TaskId,
		Retry:       req.Retry,
		MaxRetries:  req.MaxRetries,
	}

	Id, err := c.svc.Save(ctx, job, uc.Id)
	if err != nil {
		c.l.Info(fmt.Sprintf("创建 Job 失败，用户 Id：%v", uc.Id), logger.Error(err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    Id,
	}, nil
}

// Add 添加全部类型
func (c *CronJobHandler) AddAll(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	var req CronJobReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.JobInvalidInput,
			Message: "用户输入格式不正确",
		}, err
	}

	duration, _ := time.ParseDuration(req.Duration)
	job := domain.CronJob{
		Id:          req.Id, // 传入了 id 就是修改，不传 id 就是新增
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Cron:        req.Cron,
		HttpCfg:     req.HttpCfg,
		TaskId:      req.TaskId,
		TimeZone:    req.TimeZone,
		Duration:    duration,
		Retry:       req.Retry,
		MaxRetries:  req.MaxRetries,
	}

	Id, err := c.svc.Save(ctx, job, uc.Id)
	if err != nil {
		c.l.Info(fmt.Sprintf("创建 Job 失败，用户 Id：%v", uc.Id), logger.Error(err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    Id,
	}, nil
}
