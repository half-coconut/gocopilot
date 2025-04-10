package repository

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository/dao"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"database/sql"
	"time"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	Create(ctx context.Context, job domain.CronJob) (int64, error)
	Update(ctx context.Context, job domain.CronJob) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	GetJobById(ctx context.Context, jid int64) (domain.CronJob, error)
}

type CacheCronJobRepository struct {
	dao      dao.CronJobDAO
	l        logger.LoggerV1
	userRepo UserRepository
}

func NewCacheCronJobRepository(dao dao.CronJobDAO, l logger.LoggerV1, userRepo UserRepository) CronJobRepository {
	return &CacheCronJobRepository{dao: dao, l: l, userRepo: userRepo}
}

func (repo *CacheCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	job, err := repo.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	creator := repo.findUserByAPI(ctx, job)
	return repo.entityToDomain(job, creator), nil
}

func (repo *CacheCronJobRepository) GetJobById(ctx context.Context, jid int64) (domain.CronJob, error) {
	job, err := repo.dao.GetJobById(ctx, jid)
	if err != nil {
		return domain.CronJob{}, err
	}
	creator := repo.findUserByAPI(ctx, job)
	return repo.entityToDomain(job, creator), nil
}

func (c *CacheCronJobRepository) findUserByAPI(ctx context.Context, job dao.CronJob) domain.User {
	// 适合单体应用
	creator, err := c.userRepo.FindById(ctx, job.CreatorId)
	if err != nil {
		c.l.Error("查询创建人失败", logger.Error(err))
	}
	return creator
}

func (repo *CacheCronJobRepository) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return repo.dao.UpdateNextTime(ctx, id, next)
}

func (repo *CacheCronJobRepository) Create(ctx context.Context, job domain.CronJob) (int64, error) {
	return repo.dao.Insert(ctx, repo.domainToEntity(job))
}

func (repo *CacheCronJobRepository) Update(ctx context.Context, job domain.CronJob) error {
	return repo.dao.UpdateById(ctx, repo.domainToEntity(job))
}

func (repo *CacheCronJobRepository) domainToEntity(job domain.CronJob) dao.CronJob {
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

func (repo *CacheCronJobRepository) entityToDomain(job dao.CronJob, creator domain.User) domain.CronJob {
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
