package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/service/core/execution"
	"TestCopilot/TestEngine/internal/service/core/model"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type APIHandler struct {
	l       logger.LoggerV1
	svc     service.APIService
	userSvc service.UserService
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
	api.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](a.Detail))

}

type APIReq struct {
	Id      int64  `json:"id"` // 判断是否新增或修改
	Name    string `json:"name"`
	URL     string `json:"url"`
	Params  string `json:"params"`
	Type    string `json:"type"`
	Body    string `json:"body"`
	Header  string `json:"header"`
	Method  string `json:"method"`
	Project string `json:"project"`
	Debug   bool   `json:"debug"` // 是否运行 debug
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
		// TODO: 以后可能不会在 web 层处理这些内容
		report := run(req, "test@123.com")
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
			DebugResult: report,
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
		}, err
	}

}

func run(req APIReq, userEmail string) string {
	var h = make(http.Header, 0)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")

	body := []byte(`{"jsonrpc": "2.0", "method": "eth_accounts", "params": [], "id": 1}`)

	ht := model.NewHttpContent(strings.ToUpper(req.Method),
		req.URL,
		req.Params,
		body,
		h,
	)
	ws := model.WebsocketContent{}
	api := model.NewAPI(req.Name, "123", req.Type, userEmail, req.Debug, *ht, ws)

	apis := []model.API{api}
	taskConf := model.NewTaskConfig(10)
	task := model.NewTask("Debug任务 APIs API", apis, *taskConf)

	e := execution.NewExecutionLoadRun(task)

	results := make(chan []*model.HttpResult)
	var wg sync.WaitGroup
	s := &model.Subtask{
		Began: time.Now(),
	}
	wg.Add(1)
	go e.HttpRunDebug(results, &wg, s)

	go func() {
		wg.Wait()
		close(results)
	}()
	return model.FinalReport(s, results)
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
				DebugResult: src.DebugResult,
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
	detail, err := a.svc.Detail(ctx, uc.Id, req.aid)
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
		DebugResult: detail.DebugResult,
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

type APIListResponse struct {
	Interfaces []API0 `json:"interfaces"` // API 列表
	Total      int    `json:"total"`      // API 总数
}

// 前端得到的API数据
type API0 struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Params      string `json:"params"`
	Body        string `json:"body"`
	Header      string `json:"header"`
	Method      string `json:"method"`
	Type        string `json:"type"` // http/websocket
	Project     string `json:"project"`
	DebugResult string `json:"debug_result"`

	Creator string `json:"creator"`
	Updater string `json:"updater"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}

// 从JSON字符串转换回http.Header
func jsonToHeader(headerJSON string) http.Header {
	// 创建一个用于解析的map
	var headerMap map[string][]string
	err := json.Unmarshal([]byte(headerJSON), &headerMap)
	if err != nil {
		return nil
	}
	// 将map转换为http.Header
	header := make(http.Header)
	for key, values := range headerMap {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	return header
}
