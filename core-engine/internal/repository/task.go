package repository

import (
	"context"
	"database/sql"
	"github.com/ecodeclub/ekit/slice"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/cache"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) (int64, error)
	Update(ctx context.Context, task domain.Task) error
	FindByUId(ctx context.Context, uid int64) ([]domain.Task, error)
	FindByTId(ctx context.Context, tid int64) (domain.Task, error)
}

type CacheTaskRepository struct {
	dao       dao.TaskDAO
	cache     cache.TaskCache
	userCache cache.UserCache
	l         logger.LoggerV1
	userRepo  UserRepository
	apiRepo   APIRepository
}

func (repo *CacheTaskRepository) FindByTId(ctx context.Context, tid int64) (domain.Task, error) {
	select {
	case <-ctx.Done():
		return domain.Task{}, nil
	default:

		task, err := repo.cache.Get(ctx, tid)
		if err == nil {
			return task, nil
		}

		taskEntity, err := repo.dao.FindByTId(ctx, tid)

		if err != nil {
			repo.l.Error("查询任务失败：", logger.Error(err))
			return domain.Task{}, nil
		}

		var apis []int64
		apilist, err := repo.findAPIListByAIds(ctx, jsonx.JsonUnmarshal(taskEntity.AIds.String, apis))
		if err != nil {
			repo.l.Error("查询api失败", logger.Error(err))
		}
		creator, updater := repo.findUserByUId(ctx, taskEntity)
		taskDomain := repo.entityToDomain(taskEntity, creator, updater, apilist)

		err = repo.cache.Set(ctx, taskDomain)
		if err != nil {
			repo.l.Error("缓存任务失败", logger.Error(err))
		}

		return taskDomain, err
	}

}

func (repo *CacheTaskRepository) FindByUId(ctx context.Context, uid int64) ([]domain.Task, error) {
	var tasks []dao.Task
	tasks, err := repo.dao.FindByUId(ctx, uid)
	if err != nil {
		return []domain.Task{}, err
	}

	taskList := make([]domain.Task, 0)

	for _, task := range tasks {
		subTask, err := repo.FindByTId(ctx, task.Id)
		if err != nil {
			return []domain.Task{}, err
		}
		taskList = append(taskList, subTask)
	}

	return taskList, err
}

func (repo *CacheTaskRepository) findUserByUId(ctx context.Context, task dao.Task) (domain.User, domain.User) {
	select {
	case <-ctx.Done():
		return domain.User{}, domain.User{}
	default:
		cUid := task.CreatorId
		uUid := task.UpdaterId

		creator, _ := repo.userCache.Get(ctx, cUid)
		updater, err := repo.userCache.Get(ctx, uUid)
		if err == nil {
			return creator, updater
		}

		creator, err = repo.userRepo.FindById(ctx, task.CreatorId)
		if err != nil {
			repo.l.Error("查询创建人失败", logger.Error(err))
		}
		err = repo.userCache.Set(ctx, creator)
		if err != nil {
			repo.l.Error("创建创建人缓存失败", logger.Error(err))
		}

		updater, err = repo.userRepo.FindById(ctx, task.UpdaterId)
		if err != nil {
			repo.l.Error("查询更新人失败", logger.Error(err))
		}
		err = repo.userCache.Set(ctx, updater)
		if err != nil {
			repo.l.Error("创建更新人缓存失败", logger.Error(err))
		}
		return creator, updater
	}

}

func (repo *CacheTaskRepository) findAPIListByAIds(ctx context.Context, aids []int64) ([]domain.API, error) {
	var apiList []domain.API
	for _, aid := range aids {
		subApi, err := repo.apiRepo.FindByAId(ctx, aid)
		if err != nil {
			repo.l.Error("查询api失败", logger.Error(err))
			return []domain.API{}, err
		}
		apiList = append(apiList, subApi)
	}
	return apiList, nil
}

func (repo *CacheTaskRepository) Create(ctx context.Context, task domain.Task) (int64, error) {
	return repo.dao.Insert(ctx, repo.domainToEntity(task))
}

func (repo *CacheTaskRepository) Update(ctx context.Context, task domain.Task) error {
	return repo.dao.UpdateById(ctx, repo.domainToEntity(task))
}

func NewCacheTaskRepository(dao dao.TaskDAO, cache cache.TaskCache, l logger.LoggerV1, userRepo UserRepository, apiRepo APIRepository, userCache cache.UserCache) TaskRepository {
	return &CacheTaskRepository{
		dao:       dao,
		cache:     cache,
		userCache: userCache,
		l:         l,
		userRepo:  userRepo,
		apiRepo:   apiRepo,
	}
}

func (repo *CacheTaskRepository) domainToEntity(task domain.Task) dao.Task {
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

func (repo *CacheTaskRepository) entityToDomain(task dao.Task, creator, updater domain.User, apilist []domain.API) domain.Task {
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
