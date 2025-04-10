package service

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository"
	"TestCopilot/TestEngine/internal/service/core"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type CronJobService interface {
	Save(ctx context.Context, job domain.CronJob, uid int64) (int64, error)
	ResetNextTime(ctx context.Context, jid int64) error
	ExecOne(ctx *gin.Context, jid int64) error
}

type cronJobServiceImpl struct {
	l       logger.LoggerV1
	repo    repository.CronJobRepository
	taskSvc core.TaskService
}

func NewCronJobServiceImpl(l logger.LoggerV1, repo repository.CronJobRepository, taskSvc core.TaskService) CronJobService {
	return &cronJobServiceImpl{l: l, repo: repo, taskSvc: taskSvc}
}

func (svc *cronJobServiceImpl) ExecOne(ctx *gin.Context, jid int64) error {
	// 某个任务的调度器
	// 如果任务少，可以这样 一条一条独立执行
	// 如果任务多，并且是多实例部署，就需要抢占式，通过 MySQl分布式锁抢中某个任务，然后再执行
	for {
		if ctx.Err() != nil {
			// 退出调度循环
			return ctx.Err()
		}
		dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		j, err := svc.repo.GetJobById(dbCtx, jid)
		if err != nil {
			svc.l.Error("获取任务失败", logger.Error(err))
		}
		if time.Now().After(j.NextTime) && j.TaskId != 0 {
			report := svc.taskSvc.PerformanceRun(ctx, j.TaskId)

			svc.l.Info(fmt.Sprintf("定时任务执行结果：%v", report))

			err = svc.ResetNextTime(ctx, jid)
			if err != nil {
				svc.l.Error("设置下一次执行时间失败", logger.Error(err))
			}
		}
		time.Sleep(time.Second * 30)
	}
}

func (svc *cronJobServiceImpl) ResetNextTime(ctx context.Context, jid int64) error {
	j, err := svc.repo.GetJobById(ctx, jid)
	if err != nil {
		svc.l.Error("通过 jid，获取 job 失败")
	}
	next := j.SetNextTime()
	//if next.IsZero() {
	//	// 没有下一次，不用更新
	//	return svc.repo.Stop(ctx, j.Id)
	//}
	return svc.repo.UpdateNextTime(ctx, j.Id, next)
}

func (svc *cronJobServiceImpl) Save(ctx context.Context, job domain.CronJob, uid int64) (int64, error) {
	if job.Id > 0 {
		// 这里是修改
		err := svc.repo.Update(ctx, job)
		if err != nil {
			svc.l.Warn("修改失败", logger.Error(err))
		}
		return job.Id, err
	}
	// 这里是新增
	job.Creator = domain.Editor{
		Id: uid,
	}
	Id, err := svc.repo.Create(ctx, job)
	if err != nil {
		svc.l.Warn("新增失败", logger.Error(err))
	}
	return Id, err
}
