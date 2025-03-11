package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/service"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/logger"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	ug.POST("/add", ginx.WrapToken[ijwt.UserClaims](a.Add))
	ug.GET("/list", ginx.WrapToken[ijwt.UserClaims](a.List))
}

func (a *APIHandler) Add(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	type APIReq struct {
		Id      int64       `json:"id"`
		Name    string      `json:"name"`
		URL     string      `json:"url"`
		Params  string      `json:"params,omitempty"`
		Type    string      `json:"type,omitempty"`
		Body    string      `json:"data,omitempty"`
		Header  http.Header `json:"header,omitempty"`
		Method  string      `json:"method"`
		Project string      `json:"project"`
	}
	var req APIReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	api := domain.API{
		Id:      req.Id,
		Name:    req.Name,
		URL:     req.URL,
		Params:  req.Params,
		Type:    req.Type,
		Body:    req.Body,
		Header:  req.Header,
		Method:  req.Method,
		Project: req.Project,
	}
	Id, err := a.svc.Save(ctx, api, uc.Id)
	if err != nil {
		a.l.Info(fmt.Sprintf("保存笔记失败，用户 Id：%v", uc.Id), logger.Error(err))
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	return ginx.Result{
		Code:    1,
		Message: "保存成功!",
		Data:    Id,
	}, nil
}

func (a *APIHandler) List(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	type ListReq struct {
		Id int64
	}
	var req ListReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	// 要查询什么的 api
	// 按照组织下的所有 api ，-> 查询 team_id 下的所有 api
	// 按照项目进行分类，-> 查询 project_id 下的所有 api
	// 这里现按照创建人ID 进行查询
	apis, err := a.svc.List(ctx, uc.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		a.l.Info("用户校验，系统错误", logger.Error(err), logger.Int64("Id", uc.Id))
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	api0List := slice.Map[domain.API, API0](apis,
		func(idx int, src domain.API) API0 {
			return API0{
				Id:      src.Id,
				Name:    src.Name,
				URL:     src.URL,
				Params:  src.Params,
				Body:    src.Body,
				Header:  src.Header,
				Method:  src.Method,
				Type:    src.Type,
				Project: src.Project,
				Creator: src.Creator,
				Updater: src.Updater,
				Ctime:   src.Ctime.Format(time.DateTime),
				Utime:   src.Utime.Format(time.DateTime),
			}
		})

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    api0List}, nil
}

// 前端得到的API数据
type API0 struct {
	Id      int64               `json:"id"`
	Name    string              `json:"name"`
	URL     string              `json:"url"`
	Params  string              `json:"params,omitempty"`
	Body    string              `json:"data,omitempty"`
	Header  map[string][]string `json:"header,omitempty"`
	Method  string              `json:"method"`
	Type    string              `json:"type"` // http/websocket
	Project string              `json:"project"`

	Creator int64  `json:"creator"`
	Updater int64  `json:"updater"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}
