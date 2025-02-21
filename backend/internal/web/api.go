package web

import (
	"TestCopilot/backend/internal/domain"
	"TestCopilot/backend/internal/service"
	"TestCopilot/backend/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIHandler struct {
	l   logger.LoggerV1
	svc service.APIService
}

func NewAPIHandler(svc service.APIService, l logger.LoggerV1) *APIHandler {
	return &APIHandler{
		svc: svc,
		l:   l,
	}
}

func (a *APIHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/api")
	ug.POST("/add", a.Add)
	ug.GET("/list", a.List)
}

func (a *APIHandler) Add(ctx *gin.Context) {
	type APIReq struct {
		Id     int64       `json:"id"`
		Name   string      `json:"name"`
		URL    string      `json:"url"`
		Params string      `json:"params,omitempty"`
		Body   string      `json:"data,omitempty"`
		Header http.Header `json:"header,omitempty"`
		Method string      `json:"method"`
	}
	var req APIReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	cl, _ := ctx.Get("claims")
	claims, ok := cl.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		a.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}
	api := domain.API{
		Id:     req.Id,
		Name:   req.Name,
		URL:    req.URL,
		Params: req.Params,
		Body:   req.Body,
		Header: req.Header,
		Method: req.Method,
	}
	Id, err := a.svc.Save(ctx, api, claims.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		a.l.Info(fmt.Sprintf("保存笔记失败，用户 Id：%v", claims.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Message: "保存成功！", Data: Id})
}

func (a *APIHandler) List(ctx *gin.Context) {
	type ListReq struct {
		Id int64
	}
	var req ListReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	cl, _ := ctx.Get("claims")
	claims, ok := cl.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		a.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}
	// 要查询什么的 api
	// 按照组织下的所有 api ，-> 查询 team_id 下的所有 api
	// 按照项目进行分类，-> 查询 project_id 下的所有 api
	// 这里现按照创建人ID 进行查询
	api, err := a.svc.List(ctx, claims.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		a.l.Info("用户校验，系统错误", logger.Error(err), logger.Int64("Id", claims.Id))
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Message: "获取 api list 成功！", Data: api})
}
