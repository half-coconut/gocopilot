package service

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
)

type CronJobService interface {
	Save(ctx context.Context, job domain.CronJob, uid int64) (int64, error)
}

type cronJobServiceImpl struct {
	l logger.LoggerV1
}

func (svc *cronJobServiceImpl) Save(ctx context.Context, job domain.CronJob, uid int64) (int64, error) {
	panic("implement me")
}
