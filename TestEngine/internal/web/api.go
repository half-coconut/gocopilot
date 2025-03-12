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
	api := server.Group("/api")
	api.POST("/edit", ginx.WrapToken[ijwt.UserClaims](a.Edit))
	api.GET("/list", ginx.WrapToken[ijwt.UserClaims](a.List))
	api.GET("/detail:id", ginx.WrapToken[ijwt.UserClaims](a.Detail))
}

func (a *APIHandler) Edit(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	type APIReq struct {
		Id      int64  `json:"id"`
		Name    string `json:"name"`
		URL     string `json:"url"`
		Params  string `json:"params,omitempty"`
		Type    string `json:"type,omitempty"`
		Body    string `json:"body,omitempty"`
		Header  string `json:"header,omitempty"` // 注意需要校验数据类型
		Method  string `json:"method"`
		Project string `json:"project"`
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
		Id:      req.Id, // 传入了 id 就是修改，不传 id 就是新增
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

	if req.Name == "" {
		return ginx.Result{
			Code:    0,
			Message: "昵称不能为空",
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

	response := APIListResponse{
		Items: api0List,
		Total: len(apis), // apis 的长度就是 API 总数
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response}, nil
}

func (a *APIHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	// 根据 id 查询 api
	return ginx.Result{Message: "OK"}, nil
}

type APIListResponse struct {
	Items []API0 `json:"items"` // API 列表
	Total int    `json:"total"` // API 总数
}

// 前端得到的API数据
type API0 struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Params  string `json:"params,omitempty"`
	Body    string `json:"body,omitempty"`
	Header  string `json:"header,omitempty"`
	Method  string `json:"method"`
	Type    string `json:"type"` // http/websocket
	Project string `json:"project"`

	Creator int64  `json:"creator"`
	Updater int64  `json:"updater"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}
