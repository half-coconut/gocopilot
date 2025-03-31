package repository

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository/dao"
	"TestCopilot/TestEngine/pkg/jsonx"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type APIRepository interface {
	Create(ctx context.Context, api domain.API) (int64, error)
	Update(ctx context.Context, api domain.API) error
	FindByUId(ctx context.Context, uid int64) ([]domain.API, error)
	FindByAId(ctx context.Context, aid int64) (domain.API, error)
}

type CacheAPIRepository struct {
	dao      dao.APIDAO
	l        logger.LoggerV1
	userRepo UserRepository
}

func (c *CacheAPIRepository) FindByAId(ctx context.Context, aid int64) (domain.API, error) {
	api, err := c.dao.FindByAId(ctx, aid)
	if err != nil {
		return domain.API{}, nil
	}
	creator, updater := c.findUserByAPI(ctx, api)
	return c.entityToDomain(api, creator, updater), err
}

func (c *CacheAPIRepository) FindByUId(ctx context.Context, uid int64) ([]domain.API, error) {
	// 直接查库
	var api []dao.API
	api, err := c.dao.FindByUId(ctx, uid)
	if err != nil {
		return []domain.API{}, err
	}
	apiResp := make([]domain.API, 0)

	for _, a := range api {
		creator, updater := c.findUserByAPI(ctx, a)
		aResp := c.entityToDomain(a, creator, updater)
		apiResp = append(apiResp, aResp)
	}

	return apiResp, err
}

func (c *CacheAPIRepository) findUserByAPI(ctx context.Context, api dao.API) (domain.User, domain.User) {
	// 适合单体应用
	creator, err := c.userRepo.FindById(ctx, api.CreatorId)
	if err != nil {
		c.l.Error("查询创建人失败", logger.Error(err))
	}

	updater, err := c.userRepo.FindById(ctx, api.UpdaterId)
	if err != nil {
		c.l.Error("查询更新人失败", logger.Error(err))
	}
	return creator, updater
}

func NewAPIRepository(dao dao.APIDAO, l logger.LoggerV1, userRepo UserRepository) APIRepository {
	return &CacheAPIRepository{
		dao:      dao,
		l:        l,
		userRepo: userRepo,
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
	var debugRes domain.TaskDebugLog
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

func headerToJSON(header http.Header) string {
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return ""
	}
	return string(headerJSON)
}

// 从JSON字符串转换回http.Header
func jsonToHeader(headerJSON string) http.Header {
	// 创建一个用于解析的map
	var headerMap map[string][]string
	err := json.Unmarshal([]byte(headerJSON), &headerMap)
	if err != nil {
		return nil
	}
	// 将map转换为http.Header
	header := make(http.Header)
	for key, values := range headerMap {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	return header
}
