package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/errs"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/service/core/model"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/jsonx"
	"TestCopilot/TestEngine/pkg/logger"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
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
	// 创建任务,执行性能测试，支持修改参数，然后执行
	task.POST("/edit", ginx.WrapToken[ijwt.UserClaims](t.Edit))

	// 执行一次，生成性能测试报告
	task.GET("/run/once/:id", ginx.WrapToken[ijwt.UserClaims](t.RunOnce))
	// 某个任务的 debug 执行，生成接口测试报告
	task.GET("/debug/:id", ginx.WrapToken[ijwt.UserClaims](t.Debug))

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
			//duration, err = time.ParseDuration(req.Durations)
			//if err != nil {
			//	t.l.Info(fmt.Sprintf("duration 时间有误：%v", req.Durations), logger.Error(err))
			//	return ginx.Result{
			//		Code:    errs.TaskInvalidInput,
			//		Message: "Durations 输入有误",
			//	}, err
			//}
			report := t.svc.HttpRun(ctx, req.Id,
				5*time.Minute, req.Rate)

			return ginx.Result{
				Code:    1,
				Message: fmt.Sprintf("%v, OK", req.Id),
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
			Code:    0,
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
	panic("implement me")
}

// RunOnce 设置并发数等，执行一次，生成性能测试报告
func (t *TaskHandler) RunOnce(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	tid := ctx.Param("id")

	type TaskReq struct {
		tid int64 `json:"id"`
	}
	var req TaskReq

	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}
	req.tid, err = strconv.ParseInt(tid, 10, 64)
	if err != nil {
		t.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	begin := time.Now()
	results := make(chan []*model.HttpResult)
	var wg sync.WaitGroup

	wg.Add(1)
	go t.svc.HttpRunDebug(ctx, req.tid, results, &wg)

	go func() {
		wg.Wait()
		close(results)
	}()

	content := model.FinalReport(begin, results)
	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    content,
	}, nil
}

// Debug 发送请求，把任务里所有接口跑一遍
func (t *TaskHandler) Debug(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	tid := ctx.Param("id")

	type TaskReq struct {
		tid int64 `json:"id"`
	}
	var req TaskReq

	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}
	req.tid, err = strconv.ParseInt(tid, 10, 64)
	if err != nil {
		t.l.Error(fmt.Sprintf("Error converting string to int64: %v", err))
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	debug := t.svc.OnceRunDebug(ctx, req.tid)

	return ginx.Result{
		Code:    1,
		Message: fmt.Sprintf("%v, OK", tid),
		Data:    debug,
	}, nil
}
