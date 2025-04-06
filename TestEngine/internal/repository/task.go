package repository

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository/dao"
	"TestCopilot/TestEngine/pkg/jsonx"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"database/sql"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) (int64, error)
	Update(ctx context.Context, task domain.Task) error
	FindByUId(ctx context.Context, uid int64) ([]domain.Task, error)
	FindByTId(ctx context.Context, tid int64) (domain.Task, error)
}

type CacheTaskRepository struct {
	dao      dao.TaskDAO
	l        logger.LoggerV1
	userRepo UserRepository
	apiRepo  APIRepository
}

func (c CacheTaskRepository) FindByTId(ctx context.Context, tid int64) (domain.Task, error) {
	task, err := c.dao.FindByTId(ctx, tid)

	if err != nil {
		c.l.Error("FindByTId 失败：", logger.Error(err))
		return domain.Task{}, nil
	}

	var apis []int64
	apilist, err := c.findAPIListByAIds(ctx, jsonx.JsonUnmarshal(task.AIds.String, apis))
	if err != nil {
		c.l.Error("查询api失败", logger.Error(err))
	}
	creator, updater := c.findUserByUId(ctx, task)
	taskDomain := c.entityToDomain(task, creator, updater, apilist)

	return taskDomain, err
}

func (c CacheTaskRepository) FindByUId(ctx context.Context, uid int64) ([]domain.Task, error) {
	var task []dao.Task
	task, err := c.dao.FindByUId(ctx, uid)
	if err != nil {
		return []domain.Task{}, err
	}

	taskList := make([]domain.Task, 0)

	for _, t := range task {
		var apis []int64
		apilist, err := c.findAPIListByAIds(ctx, jsonx.JsonUnmarshal(t.AIds.String, apis))
		if err != nil {
			c.l.Error("查询api失败", logger.Error(err))
		}
		creator, updater := c.findUserByUId(ctx, t)
		subTask := c.entityToDomain(t, creator, updater, apilist)
		if err != nil {
			c.l.Error("查询api失败", logger.Error(err))
		}
		taskList = append(taskList, subTask)
	}

	return taskList, err
}

func (c *CacheTaskRepository) findUserByUId(ctx context.Context, task dao.Task) (domain.User, domain.User) {
	creator, err := c.userRepo.FindById(ctx, task.CreatorId)
	if err != nil {
		c.l.Error("查询创建人失败", logger.Error(err))
	}

	updater, err := c.userRepo.FindById(ctx, task.UpdaterId)
	if err != nil {
		c.l.Error("查询更新人失败", logger.Error(err))
	}
	return creator, updater
}

func (c *CacheTaskRepository) findAPIListByAIds(ctx context.Context, aids []int64) ([]domain.API, error) {
	var apiList []domain.API
	for _, aid := range aids {
		subApi, err := c.apiRepo.FindByAId(ctx, aid)
		if err != nil {
			c.l.Error("查询api失败", logger.Error(err))
			return []domain.API{}, err
		}
		apiList = append(apiList, subApi)
	}
	return apiList, nil
}

func (c CacheTaskRepository) Create(ctx context.Context, task domain.Task) (int64, error) {
	return c.dao.Insert(ctx, c.domainToEntity(task))
}

func (c CacheTaskRepository) Update(ctx context.Context, task domain.Task) error {
	return c.dao.UpdateById(ctx, c.domainToEntity(task))
}

func NewCacheTaskRepository(dao dao.TaskDAO, l logger.LoggerV1, userRepo UserRepository, apiRepo APIRepository) TaskRepository {
	return &CacheTaskRepository{
		dao:      dao,
		l:        l,
		userRepo: userRepo,
		apiRepo:  apiRepo,
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
		Rate:       task.Rate,

		CreatorId: task.Creator.Id,
		UpdaterId: task.Updater.Id,
	}
}

func (c *CacheTaskRepository) entityToDomain(task dao.Task, creator, updater domain.User, apilist []domain.API) domain.Task {
	var (
		Body   map[string]interface{}
		Header map[string]string
	)
	apis := slice.Map[domain.API, domain.TaskAPI](apilist,
		func(idx int, src domain.API) domain.TaskAPI {
			return domain.TaskAPI{
				Id:     src.Id,
				Name:   src.Name,
				URL:    src.URL,
				Params: src.Params,
				Body:   jsonx.JsonUnmarshal(src.Body, Body),
				Header: jsonx.JsonUnmarshal(src.Header, Header),
				Method: src.Method,
				Type:   src.Type,
			}
		})

	var aidsList []int64
	return domain.Task{
		Id:   task.Id,
		Name: task.Name.String,
		AIds: jsonx.JsonUnmarshal(task.AIds.String, aidsList),
		APIs: apis,

		Durations:  time.Duration(task.Durations),
		Workers:    task.Workers,
		MaxWorkers: task.MaxWorkers,
		Rate:       task.Rate,
		Creator: domain.Editor{
			Id:   creator.Id,
			Name: creator.FullName,
		},
		Updater: domain.Editor{
			Id:   updater.Id,
			Name: updater.FullName,
		},
		Ctime: time.UnixMilli(task.Ctime),
		Utime: time.UnixMilli(task.Utime),
	}
}
