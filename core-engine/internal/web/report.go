package web

import (
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	l       logger.LoggerV1
	svc     service.APIService
	userSvc service.UserService
}

func ReportAPIHandler(svc service.APIService, l logger.LoggerV1) *APIHandler {
	return &APIHandler{
		svc: svc,
		l:   l,
	}
}

func (r *ReportHandler) RegisterRoutes(server *gin.Engine) {
	api := server.Group("/api")
	api.POST("/list:taskId", r.Task)
	api.GET("/list:apiId", r.Debug)

}

func (r *ReportHandler) Task(context *gin.Context) {

}

func (r *ReportHandler) Debug(context *gin.Context) {

}
