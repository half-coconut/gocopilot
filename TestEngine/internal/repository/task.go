package repository

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository/dao"
	"TestCopilot/TestEngine/pkg/jsonx"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"database/sql"
)

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) (int64, error)
	Update(ctx context.Context, task domain.Task) error
}

type CacheTaskRepository struct {
	dao dao.TaskDAO
	l   logger.LoggerV1
}

func (c CacheTaskRepository) Create(ctx context.Context, task domain.Task) (int64, error) {
	return c.dao.Insert(ctx, c.domainToEntity(task))
}

func (c CacheTaskRepository) Update(ctx context.Context, task domain.Task) error {
	return c.dao.UpdateById(ctx, c.domainToEntity(task))
}

func NewCacheTaskRepository(dao dao.TaskDAO, l logger.LoggerV1) TaskRepository {
	return &CacheTaskRepository{
		dao: dao,
		l:   l,
	}
}

func (c *CacheTaskRepository) domainToEntity(task domain.Task) dao.Task {
	return dao.Task{
		Id: task.Id,
		Name: sql.NullString{
			String: task.Name,
			Valid:  task.Name != "",
		},
		AIds: sql.NullString{
			String: jsonx.JsonMarshal(task.AIds),
			Valid:  jsonx.JsonMarshal(task.AIds) != "",
		},

		Durations:  int64(task.Durations),
		Workers:    task.Workers,
		MaxWorkers: task.MaxWorkers,
		Timeout:    int64(task.Timeout),

		CreatorId: task.Creator.Id,
		UpdaterId: task.Updater.Id,
	}
}
