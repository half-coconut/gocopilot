package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/service/core/model"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/logger"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"time"
)

// Task 要定义是接口测试，就是所有接口请求一遍，
// 性能测试，就需要定义清楚并发数等等，这个是面向前端的接口，想清楚定义

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
	task.POST("/debug", ginx.WrapToken[ijwt.UserClaims](t.Debug))

	task.GET("/list", ginx.WrapToken[ijwt.UserClaims](t.List))
	task.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](t.Detail))
}

func (t *TaskHandler) Edit(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {

	type TaskReq struct {
		Id         int64         `json:"id"`
		Name       string        `json:"name"`
		APIs       []API0        `json:"apis"`       // 接口里可能包含 http, 也可能是 websocket
		Durations  time.Duration `json:"durations"`  // 持续时间
		Workers    uint64        `json:"workers"`    // 并发数
		MaxWorkers uint64        `json:"maxWorkers"` // 最大持续时间
		Timeout    time.Duration `json:"timeout"`    // 超时时间
	}
	var req TaskReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	api0List := slice.Map[API0, domain.APIs](req.APIs,
		func(idx int, src API0) domain.APIs {
			return domain.APIs{
				Name:   src.Name,
				URL:    src.URL,
				Params: src.Params,
				Body:   src.Body,
				Header: src.Header,
				Method: src.Method,
				Type:   src.Type,
			}
		})

	task := domain.Task{
		Id:   req.Id, // 传入了 id 就是修改，不传 id 就是新增
		Name: req.Name,
		APIs: api0List,

		Durations:  req.Durations,
		Workers:    req.Workers,
		MaxWorkers: req.MaxWorkers,
		Timeout:    req.Timeout,
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
