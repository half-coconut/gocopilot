package repository

import (
	"context"
	"database/sql"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/cache"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

//go:generate mockgen -source=api.go -package=mocks -destination=mocks/api.mock.go APIRepository
type APIRepository interface {
	Create(ctx context.Context, api domain.API) (int64, error)
	Update(ctx context.Context, api domain.API) error
	FindByUId(ctx context.Context, uid int64) ([]domain.API, error)
	FindByAId(ctx context.Context, aid int64) (domain.API, error)
}

type CacheAPIRepository struct {
	dao       dao.APIDAO
	cache     cache.APICache
	userCache cache.UserCache
	l         logger.LoggerV1
	userRepo  UserRepository
}

func (c *CacheAPIRepository) FindByAId(ctx context.Context, aid int64) (domain.API, error) {
	api, err := c.cache.Get(ctx, aid)
	if err == nil {
		return api, nil
	}

	apiEntity, err := c.dao.FindByAId(ctx, aid)
	if err != nil {
		return domain.API{}, nil
	}
	creator, updater := c.findUserByAPI(ctx, apiEntity)
	apiDomain := c.entityToDomain(apiEntity, creator, updater)

	err = c.cache.Set(ctx, apiDomain)
	if err != nil {
		c.l.Error("缓存API失败")
	}
	return apiDomain, err
}

func (c *CacheAPIRepository) FindByUId(ctx context.Context, uid int64) ([]domain.API, error) {
	// 直接查库
	var apis []dao.API
	apis, err := c.dao.FindByUId(ctx, uid)
	if err != nil {
		return []domain.API{}, err
	}
	apiList := make([]domain.API, 0)

	for _, api := range apis {
		subAPI, err := c.FindByAId(ctx, api.Id)
		if err != nil {
			return []domain.API{}, err
		}
		apiList = append(apiList, subAPI)
	}

	return apiList, err
}

func (c *CacheAPIRepository) findUserByAPI(ctx context.Context, api dao.API) (domain.User, domain.User) {
	select {
	case <-ctx.Done():
		return domain.User{}, domain.User{}
	default:
		cUid := api.CreatorId
		uUid := api.UpdaterId

		creator, err := c.userCache.Get(ctx, cUid)
		updater, err := c.userCache.Get(ctx, uUid)
		if err == nil {
			return creator, updater
		}

		creator, err = c.userRepo.FindById(ctx, api.CreatorId)
		if err != nil {
			c.l.Error("查询创建人失败", logger.Error(err))
		}
		err = c.userCache.Set(ctx, creator)
		if err != nil {
			return domain.User{}, domain.User{}
		}

		updater, err = c.userRepo.FindById(ctx, api.UpdaterId)
		if err != nil {
			c.l.Error("查询更新人失败", logger.Error(err))
		}
		err = c.userCache.Set(ctx, updater)
		if err != nil {
			return domain.User{}, domain.User{}
		}
		return creator, updater
	}
}

func NewAPIRepository(dao dao.APIDAO, cache cache.APICache, userCache cache.UserCache, l logger.LoggerV1, userRepo UserRepository) APIRepository {
	return &CacheAPIRepository{
		dao:       dao,
		cache:     cache,
		userCache: userCache,
		l:         l,
		userRepo:  userRepo,
	}
}

func (c *CacheAPIRepository) Create(ctx context.Context, api domain.API) (int64, error) {
	return c.dao.Insert(ctx, c.domainToEntity(api))
}

func (c *CacheAPIRepository) Update(ctx context.Context, api domain.API) error {
	return c.dao.UpdateById(ctx, c.domainToEntity(api))
}

func (c *CacheAPIRepository) domainToEntity(api domain.API) dao.API {
	return dao.API{
		Id: api.Id,
		Name: sql.NullString{
			String: api.Name,
			Valid:  api.Name != "",
		},
		URL: sql.NullString{
			String: api.URL,
			Valid:  api.URL != "",
		},
		Params: sql.NullString{
			String: api.Params,
			Valid:  api.Params != "",
		},
		Type: sql.NullString{
			String: api.Type,
			Valid:  api.Type != "",
		},
		Body: sql.NullString{
			String: api.Body,
			Valid:  api.Body != "",
		},
		Header: sql.NullString{
			String: api.Header,
			Valid:  api.Header != "",
		},
		Method: sql.NullString{
			String: api.Method,
			Valid:  api.Method != "",
		},
		Project: sql.NullString{
			String: api.Project,
			Valid:  api.Project != "",
		},
		DebugResult: sql.NullString{
			String: jsonx.JsonMarshal(api.DebugResult),
			Valid:  jsonx.JsonMarshal(api.DebugResult) != "",
		},
		CreatorId: api.Creator.Id,
		UpdaterId: api.Updater.Id,
	}
}

func (c *CacheAPIRepository) entityToDomain(api dao.API, creator, updater domain.User) domain.API {
	var debugRes domain.DebugLog
	return domain.API{
		Id:   api.Id,
		Name: api.Name.String,
		URL:  api.URL.String,

		Params:      api.Params.String,
		Type:        api.Type.String,
		Body:        api.Body.String,
		Header:      api.Header.String,
		Method:      api.Method.String,
		Project:     api.Project.String,
		DebugResult: jsonx.JsonUnmarshal(api.DebugResult.String, debugRes),
		Creator: domain.Editor{
			Id:   creator.Id,
			Name: creator.FullName,
		},
		Updater: domain.Editor{
			Id:   updater.Id,
			Name: updater.FullName,
		},

		Ctime: time.UnixMilli(api.Ctime),
		Utime: time.UnixMilli(api.Utime),
	}
}
