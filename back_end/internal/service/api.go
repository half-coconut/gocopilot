package service

import (
	"context"
	"egg_yolk/internal/domain"
	"egg_yolk/internal/repository"
	"egg_yolk/pkg/logger"
)

type APIService interface {
	Save(ctx context.Context, api domain.API, uid int64) (int64, error)
	List(ctx context.Context, id int64) ([]domain.API, error)
}

type apiService struct {
	repo repository.APIRepository
	l    logger.LoggerV1
}

func (a apiService) List(ctx context.Context, id int64) ([]domain.API, error) {
	return a.repo.FindByUId(ctx, id)
}

func (a apiService) Save(ctx context.Context, api domain.API, uid int64) (int64, error) {
	if api.Id > 0 {
		// 这里是修改
		api.Updater = uid
		err := a.repo.Update(ctx, api)
		if err != nil {
			a.l.Warn("修改失败", logger.Error(err))
		}
		return api.Id, err
	}
	// 这里是新增
	api.Creator = uid
	api.Updater = uid
	Id, err := a.repo.Create(ctx, api)
	if err != nil {
		a.l.Warn("新增失败", logger.Error(err))
	}
	return Id, err
}

func NewAPIService(repo repository.APIRepository, l logger.LoggerV1) APIService {
	return &apiService{
		repo: repo,
		l:    l,
	}
}
