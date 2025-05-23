package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/errs"
	"github.com/half-coconut/gocopilot/core-engine/internal/service"
	"github.com/half-coconut/gocopilot/core-engine/internal/service/core"
	ijwt "github.com/half-coconut/gocopilot/core-engine/internal/web/jwt"
	"github.com/half-coconut/gocopilot/core-engine/pkg/ginx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"strconv"
	"time"
)

var _ handler = (*CronJobHandler)(nil)

type CronJobHandler struct {
	l       logger.LoggerV1
	svc     service.CronJobService
	taskSvc core.TaskService
}

func NewCronJobHandler(l logger.LoggerV1, svc service.CronJobService, taskSvc core.TaskService) *CronJobHandler {
	return &CronJobHandler{l: l, svc: svc, taskSvc: taskSvc}
}

func (c *CronJobHandler) RegisterRoutes(server *gin.Engine) {
	job := server.Group("/job")
	// 开启就是执行
	job.GET("/open/:id", ginx.WrapToken[ijwt.UserClaims](c.Open))
	job.GET("/close/:id", ginx.WrapToken[ijwt.UserClaims](c.Close))
	job.GET("/list", ginx.WrapToken[ijwt.UserClaims](c.List))
	job.POST("/add", ginx.WrapToken[ijwt.UserClaims](c.AddAll))

	add := job.Group("/add")
	add.POST("/task", ginx.WrapToken[ijwt.UserClaims](c.AddInternalTask))
	add.POST("/http", ginx.WrapToken[ijwt.UserClaims](c.AddHttpMode))

}

func (c *CronJobHandler) Open(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	jid := ctx.Param("id")

	type JobReq struct {
		Jid int64 `json:"id"`
	}
	var req JobReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.CronJobInvalidInput,
			Message: "输入格式不正确",
		}, err
	}
	req.Jid, err = strconv.ParseInt(jid, 10, 64)
	if err != nil {
		c.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    errs.CronJobInternalServerError,
			Message: "系统错误",
		}, err
	}
	// 先释放任务
	err = c.svc.Release(ctx, req.Jid)
	if err != nil {
		return ginx.Result{
			Code:    errs.CronJobInternalServerError,
			Message: "系统错误",
		}, err
	}

	err = c.svc.ResetNextTime(ctx, req.Jid)
	if err != nil {
		return ginx.Result{
			Code:    errs.CronJobInternalServerError,
			Message: "系统错误",
		}, err
	}
	// 异步调用执行定时任务
	go func() {
		err = c.svc.ExecOne(ctx, req.Jid)
		if err != nil {
			c.l.Error("执行报错", logger.Error(err))
		}
	}()

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    jid,
	}, err
}

func (c *CronJobHandler) Close(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	jid := ctx.Param("id")

	type JobReq struct {
		Jid int64 `json:"id"`
	}
	var req JobReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.CronJobInvalidInput,
			Message: "输入格式不正确",
		}, err
	}
	req.Jid, err = strconv.ParseInt(jid, 10, 64)
	if err != nil {
		c.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    errs.CronJobInternalServerError,
			Message: "系统错误",
		}, err
	}

	err = c.svc.StopOne(ctx, req.Jid)
	if err != nil {
		c.l.Error(fmt.Sprintf("暂停任务失败: %v", err))
		return ginx.Result{
			Code:    errs.CronJobInternalServerError,
			Message: "系统错误",
		}, err
	}
	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    jid,
	}, err
}

// AddHttpMode 添加Http请求类型
func (c *CronJobHandler) AddHttpMode(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	var req CronJobReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.CronJobInvalidInput,
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
			Code:    errs.CronJobInternalServerError,
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
			Code:    errs.CronJobInvalidInput,
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
			Code:    errs.CronJobInternalServerError,
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
			Code:    errs.CronJobInvalidInput,
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
	c.l.Debug(fmt.Sprintf("打印 job 的结构体：%v", job))

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

func (c *CronJobHandler) List(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	jobs, err := c.svc.List(ctx, uc.Id)

	if err != nil {
		c.l.Info("用户校验，系统错误", logger.Error(err), logger.Int64("Id", uc.Id))
		return ginx.Result{
			Code:    errs.CronJobInternalServerError,
			Message: "系统错误",
		}, err
	}
	job0List := slice.Map[domain.CronJob, CronJob0](jobs,
		func(idx int, src domain.CronJob) CronJob0 {

			return CronJob0{
				Id:          src.Id,
				Name:        src.Name,
				Description: src.Description,
				Type:        src.Type,
				Cron:        src.Cron,
				HttpCfg:     src.HttpCfg,
				TaskId:      src.TaskId,
				TimeZone:    src.TimeZone,
				Duration:    jsonx.JsonMarshal(src.Duration),
				Retry:       src.Retry,
				MaxRetries:  src.MaxRetries,
				NextTime:    src.NextTime.Format(time.DateTime),
				Status:      src.Status,

				Creator: src.Creator.Name,
				Ctime:   src.Ctime.Format(time.DateTime),
				Utime:   src.Utime.Format(time.DateTime),
			}
		})

	response := CronJobListResponse{
		Cronjob: job0List,
		Total:   len(jobs),
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response}, nil

}
