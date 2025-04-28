package repository

import (
	"context"
	"database/sql"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/cache"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	PreemptByJId(ctx context.Context, jid int64) (domain.CronJob, error)
	Create(ctx context.Context, job domain.CronJob) (int64, error)
	Update(ctx context.Context, job domain.CronJob) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	FindByJId(ctx context.Context, jid int64) (domain.CronJob, error)
	FindByUId(ctx context.Context, uid int64) ([]domain.CronJob, error)
	GetJobStatusById(ctx context.Context, jid int64) (int, error)
	Stop(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
}

type CachedCronJobRepository struct {
	dao      dao.CronJobDAO
	cache    cache.UserCache
	l        logger.LoggerV1
	userRepo UserRepository
}

func NewCacheCronJobRepository(dao dao.CronJobDAO, cache cache.UserCache, l logger.LoggerV1, userRepo UserRepository) CronJobRepository {
	return &CachedCronJobRepository{dao: dao, cache: cache, l: l, userRepo: userRepo}
}

func (repo *CachedCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	job, err := repo.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	creator := repo.findUserByAPI(ctx, job)
	return repo.entityToDomain(job, creator), nil
}

func (repo *CachedCronJobRepository) PreemptByJId(ctx context.Context, jid int64) (domain.CronJob, error) {
	job, err := repo.dao.PreemptByJId(ctx, jid)
	if err != nil {
		return domain.CronJob{}, err
	}
	creator := repo.findUserByAPI(ctx, job)
	return repo.entityToDomain(job, creator), nil
}

func (repo *CachedCronJobRepository) GetJobStatusById(ctx context.Context, jid int64) (int, error) {
	return repo.dao.GetJobStatusById(ctx, jid)
}

func (repo *CachedCronJobRepository) FindByUId(ctx context.Context, uid int64) ([]domain.CronJob, error) {
	jobs, err := repo.dao.GetJobByUId(ctx, uid)
	if err != nil {
		return []domain.CronJob{}, err
	}
	jobList := make([]domain.CronJob, 0)
	for _, job := range jobs {
		subJob, err := repo.FindByJId(ctx, job.Id)
		if err != nil {
			return []domain.CronJob{}, err
		}
		jobList = append(jobList, subJob)
	}
	return jobList, nil
}

func (repo *CachedCronJobRepository) FindByJId(ctx context.Context, jid int64) (domain.CronJob, error) {
	job, err := repo.dao.GetJobByJId(ctx, jid)
	if err != nil {
		return domain.CronJob{}, err
	}
	creator := repo.findUserByAPI(ctx, job)
	return repo.entityToDomain(job, creator), nil
}

func (c *CachedCronJobRepository) findUserByAPI(ctx context.Context, job dao.CronJob) domain.User {
	cUid := job.CreatorId

	creator, err := c.cache.Get(ctx, cUid)
	if err == nil {
		return creator
	}

	creator, err = c.userRepo.FindById(ctx, job.CreatorId)
	if err != nil {
		c.l.Error("查询创建人失败", logger.Error(err))
	}
	err = c.cache.Set(ctx, creator)
	if err != nil {
		return domain.User{}
	}

	return creator
}

func (repo *CachedCronJobRepository) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return repo.dao.UpdateNextTime(ctx, id, next)
}

func (repo *CachedCronJobRepository) Stop(ctx context.Context, id int64) error {
	return repo.dao.Stop(ctx, id)
}

func (repo *CachedCronJobRepository) Release(ctx context.Context, id int64) error {
	return repo.dao.Release(ctx, id)
}

func (repo *CachedCronJobRepository) Create(ctx context.Context, job domain.CronJob) (int64, error) {
	return repo.dao.Insert(ctx, repo.domainToEntity(job))
}

func (repo *CachedCronJobRepository) Update(ctx context.Context, job domain.CronJob) error {
	return repo.dao.UpdateById(ctx, repo.domainToEntity(job))
}

func (repo *CachedCronJobRepository) domainToEntity(job domain.CronJob) dao.CronJob {
	return dao.CronJob{
		Id: job.Id,
		Name: sql.NullString{
			String: job.Name,
			Valid:  job.Name != "",
		},
		Description: sql.NullString{
			String: job.Description,
			Valid:  job.Description != "",
		},
		Type: sql.NullString{
			String: job.Type,
			Valid:  job.Type != "",
		},
		Cron: sql.NullString{
			String: job.Cron,
			Valid:  job.Cron != "",
		},
		HttpCfg: sql.NullString{
			String: job.HttpCfg,
			Valid:  job.HttpCfg != "",
		},

		TaskId:     job.TaskId,
		Duration:   int64(job.Duration),
		Retry:      job.Retry,
		MaxRetries: job.MaxRetries,
		NextTime:   job.NextTime.UnixMilli(),

		CreatorId: job.Creator.Id,
	}
}

func (repo *CachedCronJobRepository) entityToDomain(job dao.CronJob, creator domain.User) domain.CronJob {
	return domain.CronJob{
		Id:          job.Id,
		Name:        job.Name.String,
		Description: job.Description.String,
		Type:        job.Type.String,
		Cron:        job.Cron.String,
		HttpCfg:     job.HttpCfg.String,

		TaskId:     job.TaskId,
		Duration:   time.Duration(job.Duration),
		Retry:      job.Retry,
		MaxRetries: job.MaxRetries,
		NextTime:   time.UnixMilli(job.NextTime),

		Creator: domain.Editor{
			Id:   creator.Id,
			Name: creator.FullName,
		},

		Ctime: time.UnixMilli(job.Ctime),
		Utime: time.UnixMilli(job.Utime),
	}
}
