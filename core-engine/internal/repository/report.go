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
	FindByUId(ctx context.Context, uid int64) ([]domain.Task, error)
}

type UncachedReportRepository struct {
	dao dao.ReportDAO
	l   logger.LoggerV1
}

func NewUncachedReportRepository(dao dao.ReportDAO, l logger.LoggerV1) ReportRepository {
	return &UncachedReportRepository{dao: dao, l: l}
}

func (c *UncachedReportRepository) CreateSummary(ctx context.Context, s domain.Summary) (int64, error) {
	return c.dao.InsertSummary(ctx, c.ToEntity(s))
}

func (c *UncachedReportRepository) CreateDebugLog(ctx context.Context, logs domain.DebugLog) (int64, error) {
	return c.dao.InsertDebugLog(ctx, c.domainToEntity(logs))
}

func (c *UncachedReportRepository) FindByTId(ctx context.Context, tid int64) (domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (c *UncachedReportRepository) FindByUId(ctx context.Context, uid int64) ([]domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (c *UncachedReportRepository) domainToEntity(logs domain.DebugLog) dao.DebugLog {
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

func (c *UncachedReportRepository) ToEntity(s domain.Summary) dao.Summary {
	return dao.Summary{
		Id:            s.Id,
		TaskId:        s.TaskId,
		AIds:          s.AIds,
		TName:         s.TName,
		Debug:         s.Debug,
		Total:         s.Total,
		Rate:          s.Rate,
		Throughput:    s.Throughput,
		TotalDuration: int64(s.TotalDuration),
		Min:           int64(s.Min),
		Mean:          int64(s.Mean),
		Max:           int64(s.Max),
		P50:           int64(s.P50),
		P90:           int64(s.P90),
		P95:           int64(s.P95),
		P99:           int64(s.P99),
		Ratio:         s.Ratio,
		TestStatus:    dao.TestStatus(s.TestStatus),
		StatusCodes:   s.StatusCodes,

		Status: s.Status,
	}
}
