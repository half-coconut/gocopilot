package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/service/core/model"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/jsonx"
	"TestCopilot/TestEngine/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// Task 接口测试，性能测试的任务
// 性能测试，就需要定义清楚并发数，done

type TaskHandler struct {
	l       logger.LoggerV1
	svc     model.TaskService
	userSvc service.UserService
}

func NewTaskHandler(l logger.LoggerV1, svc model.TaskService, userSvc service.UserService) *TaskHandler {
	return &TaskHandler{
		l:       l,
		svc:     svc,
		userSvc: userSvc}
}

func (t *TaskHandler) RegisterRoutes(server *gin.Engine) {
	task := server.Group("/task")
	// 创建任务，并执行
	task.POST("/edit", ginx.WrapToken[ijwt.UserClaims](t.Edit))
	task.POST("/run", ginx.WrapToken[ijwt.UserClaims](t.Run))
	// 某个任务的 debug 执行
	task.GET("/debug/:id", ginx.WrapToken[ijwt.UserClaims](t.Debug))

	task.GET("/list", ginx.WrapToken[ijwt.UserClaims](t.List))
	task.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](t.Detail))
}

func (t *TaskHandler) Edit(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	var req TaskReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	var Aids []int64
	duration, _ := time.ParseDuration(req.Durations)
	timeout, _ := time.ParseDuration(req.Timeout)

	task := domain.Task{
		Id:   req.Id, // 传入了 id 就是修改，不传 id 就是新增
		Name: req.Name,
		AIds: jsonx.JsonUnmarshal(req.AIds, Aids),

		Durations:  duration,
		Workers:    uint64(req.Workers),
		MaxWorkers: uint64(req.MaxWorkers),
		Timeout:    timeout,
	}
	Id, err := t.svc.Save(ctx, task, uc.Id)
	if err != nil {
		t.l.Info(fmt.Sprintf("创建任务失败，用户 Id：%v", uc.Id), logger.Error(err))
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}
	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    Id,
	}, nil
}

func (t *TaskHandler) List(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	panic("implement me")
}

func (t *TaskHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	panic("implement me")
}

func (t *TaskHandler) Run(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	panic("implement me")
}

func (t *TaskHandler) Debug(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	panic("implement me")
}
