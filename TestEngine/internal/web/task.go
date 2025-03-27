package web

import (
	"TestCopilot/TestEngine/internal/service"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/logger"
	"github.com/gin-gonic/gin"
)

// task 里包含调用 api 和 case
// case
// http 请求的任务，+ websocket 请求的任务 --> service 层
type TaskHandler struct {
	l logger.LoggerV1
	//svc     service.TaskService
	userSvc service.UserService
}

func (t *TaskHandler) RegisterRoutes(server *gin.Engine) {
	task := server.Group("/task")
	task.POST("/edit", ginx.WrapToken[ijwt.UserClaims](t.Edit))
	task.GET("/list", ginx.WrapToken[ijwt.UserClaims](t.List))
	task.GET("/detail:id", ginx.WrapToken[ijwt.UserClaims](t.Detail))
}

func (t *TaskHandler) Edit(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	panic("implement me")
}

func (t *TaskHandler) List(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	panic("implement me")
}

func (t *TaskHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	panic("implement me")
}
