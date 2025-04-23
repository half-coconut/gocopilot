package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/errs"
	"github.com/half-coconut/gocopilot/core-engine/internal/service/core"
	ijwt "github.com/half-coconut/gocopilot/core-engine/internal/web/jwt"
	"github.com/half-coconut/gocopilot/core-engine/pkg/ginx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"strconv"
	"sync"
	"time"
)

var _ handler = (*TaskHandler)(nil)

// Task 接口测试，性能测试的任务

type TaskHandler struct {
	l         logger.LoggerV1
	svc       core.TaskService
	reportSvc core.ReportService
}

func NewTaskHandler(l logger.LoggerV1, svc core.TaskService, reportSvc core.ReportService) *TaskHandler {
	return &TaskHandler{
		l:         l,
		svc:       svc,
		reportSvc: reportSvc,
	}
}

func (t *TaskHandler) RegisterRoutes(server *gin.Engine) {
	task := server.Group("/task")
	// 创建任务
	// 设置或者修改 rate 和 duration，并执行性能测试
	task.POST("/edit", ginx.WrapToken[ijwt.UserClaims](t.Edit))

	// Execute
	// 执行性能测试， 按照持续时间和 rate limit 来执行
	task.GET("/execute/:id", ginx.WrapToken[ijwt.UserClaims](t.Execute))
	// Debug
	// 执行一次性能测试，设置 rate 和 duration ，生成性能测试报告
	task.GET("/debug/:id", ginx.WrapToken[ijwt.UserClaims](t.PerformanceDebug))
	// 执行一次接口测试，某个任务的接口调试，生成接口测试报告
	task.GET("/debug/interfaces/:id", ginx.WrapToken[ijwt.UserClaims](t.InterfaceDebug))

	task.GET("/list", ginx.WrapToken[ijwt.UserClaims](t.List))
	task.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](t.Detail))
}

func (t *TaskHandler) Edit(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	var req TaskReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.TaskInvalidInput,
			Message: "用户输入格式不正确",
		}, err
	}

	duration, _ := time.ParseDuration(req.Durations)

	task := domain.Task{
		Id:   req.Id, // 传入了 id 就是修改，不传 id 就是新增
		Name: req.Name,
		AIds: req.AIds,

		Durations:  duration,
		Workers:    uint64(req.Workers),
		MaxWorkers: uint64(req.MaxWorkers),
		Rate:       float64(req.Rate),
	}
	Id, err := t.svc.Save(ctx, task, uc.Id)
	if err != nil {
		t.l.Info(fmt.Sprintf("创建任务失败，用户 Id：%v", uc.Id), logger.Error(err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	if !req.Execute {
		return ginx.Result{
			Code:    1,
			Message: "OK",
			Data:    Id,
		}, nil
	} else {
		if (req.Id > 0) && (req.Execute) {

			report := t.svc.ExecutePerformanceTask(ctx, req.Id)

			return ginx.Result{
				Code:    1,
				Message: fmt.Sprintf("OK"),
				Data:    report,
			}, nil
		}
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
}

func (t *TaskHandler) List(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	type ListReq struct {
		Id int64
	}
	var req ListReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}

	tasks, err := t.svc.List(ctx, uc.Id)

	if err != nil {
		t.l.Info("用户校验，系统错误", logger.Error(err), logger.Int64("Id", uc.Id))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}

	task0List := slice.Map[domain.Task, Task0](tasks,
		func(idx int, src domain.Task) Task0 {

			var aidsList []string
			for _, id := range src.AIds {
				aidsList = append(aidsList, strconv.FormatInt(id, 10))
			}
			var apiNameList []string
			for _, api := range src.APIs {
				apiNameList = append(apiNameList, api.Name)
			}
			return Task0{
				Id:         src.Id,
				Name:       src.Name,
				AIds:       aidsList, // 把接口 Name 返给前端
				APIs:       apiNameList,
				Durations:  jsonx.JsonMarshal(src.Durations),
				Workers:    src.Workers,
				MaxWorkers: src.MaxWorkers,
				Rate:       src.Rate,
				Creator:    src.Creator.Name,
				Updater:    src.Updater.Name,
				Ctime:      src.Ctime.Format(time.DateTime),
				Utime:      src.Utime.Format(time.DateTime),
			}
		})

	response := TaskListResponse{
		Tasks: task0List,
		Total: len(tasks),
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response}, nil
}

func (t *TaskHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	tid := ctx.Param("id")

	type TaskReq struct {
		tid int64 `json:"id"`
	}
	var req TaskReq

	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	req.tid, err = strconv.ParseInt(tid, 10, 64)
	if err != nil {
		t.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}

	tasks, err := t.svc.GetDetailByTid(ctx, req.tid)

	if err != nil {
		t.l.Info("用户校验，系统错误", logger.Error(err), logger.Int64("Id", uc.Id))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}

	var aidsList []string
	for _, id := range tasks.AIds {
		aidsList = append(aidsList, strconv.FormatInt(id, 10))
	}
	var apiNameList []string
	for _, api := range tasks.APIs {
		apiNameList = append(apiNameList, api.Name)
	}
	response := Task0{
		Id:         tasks.Id,
		Name:       tasks.Name,
		AIds:       aidsList, // 把接口 Name 返给前端
		APIs:       apiNameList,
		Durations:  jsonx.JsonMarshal(tasks.Durations),
		Workers:    tasks.Workers,
		MaxWorkers: tasks.MaxWorkers,
		Rate:       tasks.Rate,
		Creator:    tasks.Creator.Name,
		Updater:    tasks.Updater.Name,
		Ctime:      tasks.Ctime.Format(time.DateTime),
		Utime:      tasks.Utime.Format(time.DateTime),
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response}, nil
}

func (t *TaskHandler) Execute(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	tid := ctx.Param("id")

	type TaskReq struct {
		tid int64 `json:"id"`
	}
	var req TaskReq

	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	req.tid, err = strconv.ParseInt(tid, 10, 64)
	if err != nil {
		t.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}

	report := t.svc.ExecutePerformanceTask(ctx, req.tid)

	return ginx.Result{
		Code:    1,
		Message: fmt.Sprintf("OK"),
		Data:    report,
	}, nil
}

// PerformanceDebug 设置并发数等，执行一次，生成性能测试报告
func (t *TaskHandler) PerformanceDebug(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	tid := ctx.Param("id")

	type TaskReq struct {
		tid int64 `json:"id"`
	}
	var req TaskReq

	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	req.tid, err = strconv.ParseInt(tid, 10, 64)
	if err != nil {
		t.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}

	t.svc.SetBegin(ctx)
	begin := time.Now()
	results := make(chan []*domain.HttpResult)
	var wg sync.WaitGroup

	wg.Add(1)
	go t.svc.RunPerformanceWithDebug(ctx, req.tid, results, &wg)

	go func() {
		wg.Wait()
		close(results)
	}()

	content := t.reportSvc.GenerateReport(begin, results)
	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    content,
	}, nil
}

// InterfaceDebug 发送请求，把任务里所有接口跑一遍
func (t *TaskHandler) InterfaceDebug(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	tid := ctx.Param("id")

	type TaskReq struct {
		tid int64 `json:"id"`
	}
	var req TaskReq

	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}
	req.tid, err = strconv.ParseInt(tid, 10, 64)
	if err != nil {
		t.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    errs.TaskInternalServerError,
			Message: "系统错误",
		}, err
	}

	debug := t.svc.GetAPIDebugLogs(ctx, req.tid)

	return ginx.Result{
		Code:    1,
		Message: fmt.Sprintf("%v, OK", tid),
		Data:    debug,
	}, nil
}
