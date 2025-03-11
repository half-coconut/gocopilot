package repository

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository/dao"
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
	FindByUId(ctx context.Context, id int64) ([]domain.API, error)
}

type CacheAPIRepository struct {
	dao dao.APIDAO
	l   logger.LoggerV1
}

func (c *CacheAPIRepository) FindByUId(ctx context.Context, id int64) ([]domain.API, error) {
	// 直接查库
	var api []dao.API
	api, err := c.dao.FindByUId(ctx, id)
	if err != nil {
		return []domain.API{}, err
	}
	apiResp := make([]domain.API, 0)

	for _, a := range api {
		aResp := c.entityToDomain(a)
		apiResp = append(apiResp, aResp)
	}

	return apiResp, err
}

func NewAPIRepository(dao dao.APIDAO, l logger.LoggerV1) APIRepository {
	return &CacheAPIRepository{
		dao: dao,
		l:   l,
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
			String: headerToJSON(api.Header),
			Valid:  headerToJSON(api.Header) != "",
		},
		Method: sql.NullString{
			String: api.Method,
			Valid:  api.Method != "",
		},
		Project: sql.NullString{
			String: api.Project,
			Valid:  api.Project != "",
		},
		Creator: api.Creator,
		Updater: api.Updater,
	}
}

func (c *CacheAPIRepository) entityToDomain(api dao.API) domain.API {
	return domain.API{
		Id:   api.Id,
		Name: api.Name.String,
		URL:  api.URL.String,

		Params:  api.Params.String,
		Type:    api.Type.String,
		Body:    api.Body.String,
		Header:  jsonToHeader(api.Header.String),
		Method:  api.Method.String,
		Project: api.Project.String,
		Creator: api.Creator,
		Updater: api.Updater,

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
