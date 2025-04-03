package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/service/core"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/jsonx"
	"TestCopilot/TestEngine/pkg/logger"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type APIHandler struct {
	l       logger.LoggerV1
	svc     service.APIService
	taskSvc core.TaskService
	userSvc service.UserService
}

func NewAPIHandler(svc service.APIService, taskSvc core.TaskService, userSvc service.UserService, l logger.LoggerV1) *APIHandler {
	return &APIHandler{
		svc:     svc,
		taskSvc: taskSvc,
		userSvc: userSvc,
		l:       l,
	}
}

func (a *APIHandler) RegisterRoutes(server *gin.Engine) {
	api := server.Group("/api")
	api.POST("/edit", ginx.WrapToken[ijwt.UserClaims](a.Edit))
	api.GET("/list", ginx.WrapToken[ijwt.UserClaims](a.List))
	api.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](a.Detail))

}

func (a *APIHandler) Edit(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	// 新增，修改 和 debug 功能
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
		Method:  strings.ToUpper(req.Method),
		Project: req.Project,
		Debug:   req.Debug,
	}

	if req.Name == "" {
		return ginx.Result{
			Code:    0,
			Message: "名称不能为空",
		}, err
	}

	// 只要修改了，就先保存
	Id, err := a.svc.Save(ctx, api, uc.Id)
	if err != nil {
		a.l.Info(fmt.Sprintf("保存笔记失败，用户 Id：%v", uc.Id), logger.Error(err))
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	if !req.Debug {
		return ginx.Result{
			Code:    1,
			Message: "保存成功!",
			Data:    Id,
		}, nil
	} else {

		now := time.Now().UnixMilli()
		// 注意这里将 domain.API 转为 domain.TaskAPI
		apis := &domain.TaskAPI{
			Id:     api.Id,
			Name:   api.Name,
			URL:    api.URL,
			Params: api.Params,
			Type:   api.Type,
			Body:   jsonx.JsonUnmarshal(api.Body, Body),
			Header: jsonx.JsonUnmarshal(api.Header, Header),
			Method: api.Method,
		}

		apiList := make([]domain.TaskAPI, 0)
		apiList = append(apiList, *apis)
		task := domain.Task{
			Id:         1,
			Name:       "system_debugging",
			APIs:       apiList,
			AIds:       []int64{1},
			Durations:  DefaultDurations,
			Workers:    DefaultWorkers,
			MaxWorkers: DefaultMaxWorkers,
			Rate:       DefaultRate,
			Creator: domain.Editor{
				Id:   1,
				Name: "Egg Yolk",
			},
			Updater: domain.Editor{
				Id:   1,
				Name: "Egg Yolk",
			},
			Ctime: time.UnixMilli(now),
			Utime: time.UnixMilli(now),
		}
		display(task)
		report := a.taskSvc.DebugForAPI(ctx, task)

		// 把 debug 结果写入数据库
		api = domain.API{
			Id:          req.Id, // 传入了 id 就是修改，不传 id 就是新增
			Name:        req.Name,
			URL:         req.URL,
			Params:      req.Params,
			Type:        req.Type,
			Body:        req.Body,
			Header:      req.Header,
			Method:      strings.ToUpper(req.Method),
			Project:     req.Project,
			Debug:       req.Debug,
			DebugResult: domain.TaskDebugLog(report),
		}

		Id, err = a.svc.Save(ctx, api, uc.Id)
		if err != nil {
			a.l.Info(fmt.Sprintf("保存笔记失败，用户 Id：%v", uc.Id), logger.Error(err))
			return ginx.Result{
				Code:    0,
				Message: "系统错误",
			}, err
		}

		return ginx.Result{
			Code:    1,
			Message: fmt.Sprintf("%d, OK", Id),
			Data:    report,
		}, nil
	}
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
				Id:          src.Id,
				Name:        src.Name,
				URL:         src.URL,
				Params:      src.Params,
				Body:        src.Body,
				Header:      src.Header,
				Method:      src.Method,
				Type:        src.Type,
				Project:     src.Project,
				DebugResult: jsonx.JsonMarshal(src.DebugResult),
				Creator:     src.Creator.Name,
				Updater:     src.Updater.Name,
				Ctime:       src.Ctime.Format(time.DateTime),
				Utime:       src.Utime.Format(time.DateTime),
			}
		})

	response := APIListResponse{
		Interfaces: api0List,
		Total:      len(apis),
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response}, nil
}

func (a *APIHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	aid := ctx.Param("id")

	type APIReq struct {
		aid int64 `json:"id"`
	}
	var req APIReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}
	req.aid, err = strconv.ParseInt(aid, 10, 64)
	if err != nil {
		a.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}
	detail, err := a.svc.Detail(ctx, req.aid)
	if err != nil {
		return ginx.Result{}, err
	}

	response := API0{
		Id:          detail.Id,
		Name:        detail.Name,
		URL:         detail.URL,
		Params:      detail.Params,
		Body:        detail.Body,
		Header:      detail.Header,
		Method:      detail.Method,
		Type:        detail.Type,
		Project:     detail.Project,
		DebugResult: jsonx.JsonMarshal(detail.DebugResult),
		Creator:     detail.Creator.Name,
		Updater:     detail.Updater.Name,
		Ctime:       detail.Ctime.Format(time.DateTime),
		Utime:       detail.Utime.Format(time.DateTime),
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response,
	}, err

}

func display(task domain.Task) string {

	content := fmt.Sprintf(`
+++++ task InterfacesDebug Log: +++++
[Id: %v]
[Name: %v]
[TaskAPI: %v]
[TaskApiIds: %v]
[Durations: %v]
[Workers:%v]
[MaxWorkers: %v]
[Rate: %v]
[Creator: %v]
[Updater: %v]
[Ctime: %v]
[Utime: %v]`,
		task.Id,
		task.Name,
		task.APIs,
		task.AIds,
		task.Durations,
		task.Workers,
		task.MaxWorkers,
		task.Rate,
		task.Creator,
		task.Updater,
		task.Ctime,
		task.Utime)
	log.Println(content)
	return content

}
