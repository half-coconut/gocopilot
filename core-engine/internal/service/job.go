package service

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, j domain.Job) error
	// 返回一个释放的方法，然后调用者去调用
}

func (p *cronJobService) ResetNextTime(ctx context.Context, j domain.Job) error {
	next := j.NextTime()
	if next.IsZero() {
		// 没有下一次，不用更新
		return p.repo.Stop(ctx, j.Id)
	}
	return p.repo.UpdateNextTime(ctx, j.Id, next)
}

type cronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logger.LoggerV1
}

func newCronJobService(repo repository.JobRepository, refreshInterval time.Duration, l logger.LoggerV1) JobService {
	return &cronJobService{repo: repo, refreshInterval: refreshInterval, l: l}
}

func (p *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := p.repo.Preempt(ctx)

	//ch := make(chan struct{})
	//go func() {
	//	ticker := time.NewTicker(p.refreshInterval)
	//	for {
	//		select {
	//		case <-ticker.C:
	//			p.Refresh(j.Id)
	//		case <-ch:
	//			return
	//		}
	//	}
	//}()

	ticker := time.NewTicker(p.refreshInterval)
	go func() {
		for range ticker.C {
			p.Refresh(j.Id)
		}
	}()

	j.CancelFunc = func() error {
		//close(ch)
		ticker.Stop()
		// 在这里释放
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		return p.repo.Release(ctx, j.Id)
	}
	return j, err
}

func (p *cronJobService) Refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 续约，更新 utime
	// 处于 running，但是更新时间在 3min 之前
	err := p.repo.UpdateUtime(ctx, id)
	if err != nil {
		// 可以重试
		p.l.Error("续约失败", logger.Error(err), logger.Int64("jid", id))
	}
}
