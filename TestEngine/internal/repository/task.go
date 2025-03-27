package repository

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository/dao"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
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
	//TODO implement me
	panic("implement me")
}

func (c CacheTaskRepository) Update(ctx context.Context, task domain.Task) error {
	//TODO implement me
	panic("implement me")
}

func NewCacheTaskRepository(dao dao.TaskDAO, l logger.LoggerV1) TaskRepository {
	return &CacheTaskRepository{
		dao: dao,
		l:   l,
	}
}
