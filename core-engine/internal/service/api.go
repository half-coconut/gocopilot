package service

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
)

type APIService interface {
	Save(ctx context.Context, api domain.API, uid int64) (int64, error)
	List(ctx context.Context, uid int64) ([]domain.API, error)
	Detail(ctx context.Context, aid int64) (domain.API, error)
}

type apiService struct {
	repo repository.APIRepository
	l    logger.LoggerV1
}

func (a apiService) Detail(ctx context.Context, aid int64) (domain.API, error) {
	return a.repo.FindByAId(ctx, aid)
}

func (a apiService) List(ctx context.Context, uid int64) ([]domain.API, error) {
	return a.repo.FindByUId(ctx, uid)
}

func (a apiService) Save(ctx context.Context, api domain.API, uid int64) (int64, error) {
	if api.Id > 0 {
		// 这里是修改
		api.Updater = domain.Editor{
			Id: uid,
		}
		err := a.repo.Update(ctx, api)
		if err != nil {
			a.l.Warn("修改失败", logger.Error(err))
		}
		return api.Id, err
	}
	// 这里是新增
	api.Creator = domain.Editor{
		Id: uid,
	}
	api.Updater = domain.Editor{
		Id: uid,
	}
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
