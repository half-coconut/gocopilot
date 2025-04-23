package repository

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/pkg/jsonx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
)

type ReportRepository interface {
	CreateDebugLog(ctx context.Context, logs domain.DebugLog) (int64, error)
	CreateSummary(ctx context.Context, s domain.Summary) (int64, error)
	FindByTId(ctx context.Context, tid int64) (domain.Task, error)
}

type CacheReportRepository struct {
	dao dao.ReportDAO
	l   logger.LoggerV1
}

func NewCacheReportRepository(dao dao.ReportDAO, l logger.LoggerV1) ReportRepository {
	return &CacheReportRepository{dao: dao, l: l}
}

func (c *CacheReportRepository) CreateSummary(ctx context.Context, s domain.Summary) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CacheReportRepository) FindByTId(ctx context.Context, tid int64) (domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CacheReportRepository) CreateDebugLog(ctx context.Context, logs domain.DebugLog) (int64, error) {
	return c.dao.InsertDebugLog(ctx, c.domainToEntity(logs))

}

func (c *CacheReportRepository) domainToEntity(logs domain.DebugLog) dao.DebugLog {
	var h map[string]string
	return dao.DebugLog{
		Id:     logs.Id,
		TaskId: logs.TaskId,
		AId:    logs.AId,
		AName:  logs.AName,
		Request: dao.RequestInfo{
			URL:     logs.Url,
			Method:  logs.Method,
			Headers: jsonx.JsonUnmarshal(logs.Headers, h),
			Body:    logs.Body,
		},
		Response: logs.Response,
		ClientIP: logs.ClientIP,
		Error:    logs.Error,
	}
}
