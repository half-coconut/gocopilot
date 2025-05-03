package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/internal/service/core"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

type CronJobService interface {
	Save(ctx context.Context, job domain.CronJob, uid int64) (int64, error)
	ResetNextTime(ctx context.Context, jid int64) error
	Release(ctx context.Context, jid int64) error
	ExecOne(ctx *gin.Context, jid int64) error
	StopOne(ctx *gin.Context, jid int64) error

	List(ctx context.Context, uid int64) ([]domain.CronJob, error)
}

type CronJobServiceImpl struct {
	l        logger.LoggerV1
	repo     repository.CronJobRepository
	taskSvc  core.TaskService
	interval time.Duration
	//limiter  *semaphore.Weighted
}

func NewCronJobService(l logger.LoggerV1, repo repository.CronJobRepository, taskSvc core.TaskService) CronJobService {
	return &CronJobServiceImpl{
		l:        l,
		repo:     repo,
		taskSvc:  taskSvc,
		interval: time.Second * 10,
		//limiter:  semaphore.NewWeighted(200),
	}
}

const (
	// 等待，准备进入
	//cronjobStatusWaiting int = 0
	// 执行中
	//cronjobStatusRunning int = 1
	// 暂停调度
	cronjobStatusPaused int = 2
)

func (svc *CronJobServiceImpl) List(ctx context.Context, uid int64) ([]domain.CronJob, error) {
	return svc.repo.FindByUId(ctx, uid)
}

func (svc *CronJobServiceImpl) StopOne(ctx *gin.Context, jid int64) error {
	return svc.repo.Stop(ctx, jid)
}

func (svc *CronJobServiceImpl) ExecOne(ctx *gin.Context, jid int64) error {
	// 使用一个goroutine 执行某个任务
	// 如果任务少，可以这样 一条一条独立执行
	// 如果任务多，并且是多实例部署，就需要抢占式，通过 MySQl分布式锁抢中某个任务，然后再执行

	for {
		if ctx.Err() != nil {
			// 退出调度循环
			return ctx.Err()
		}
		//err := svc.limiter.Acquire(ctx, 1)
		//if err != nil {
		//	return err
		//}

		dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		status, err := svc.repo.GetJobStatusById(dbCtx, jid)
		if status == cronjobStatusPaused {
			return errors.New("任务已暂停")
		}

		j, err := svc.repo.PreemptByJId(dbCtx, jid)
		if err != nil {
			svc.l.Error("获取任务失败", logger.Error(err))
		}
		go func() {
			defer func() {
				//svc.limiter.Release(1)
				err = svc.repo.Release(ctx, jid)
				if err != nil {
					svc.l.Error("释放 job 失败", logger.Error(err))
				}
				err = svc.ResetNextTime(ctx, jid)
				if err != nil {
					svc.l.Error("设置下一次执行时间失败", logger.Error(err))
				}
			}()
			if j.TaskId != 0 {
				report := svc.taskSvc.ExecutePerformanceTask(ctx, j.TaskId, false)
				svc.l.Info(fmt.Sprintf("定时任务执行结果：%v", report))
			}
		}()

		time.Sleep(svc.interval)
	}
}

func (svc *CronJobServiceImpl) Release(ctx context.Context, jid int64) error {
	return svc.repo.Release(ctx, jid)
}

func (svc *CronJobServiceImpl) ResetNextTime(ctx context.Context, jid int64) error {
	j, err := svc.repo.FindByJId(ctx, jid)
	if err != nil {
		svc.l.Error("通过 jid，获取 job 失败")
	}
	next := j.SetNextTime()
	if next.IsZero() {
		// 没有下一次，不用更新
		return svc.repo.Stop(ctx, j.Id)
	}
	return svc.repo.UpdateNextTime(ctx, j.Id, next)
}

func (svc *CronJobServiceImpl) Save(ctx context.Context, job domain.CronJob, uid int64) (int64, error) {
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
