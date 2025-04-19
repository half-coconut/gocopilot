package ioc

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/job"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"time"
)

func InitScheduler(l logger.LoggerV1, local *job.LocalFuncExecutor, svc service.JobService) *job.Schedule {
	res := job.NewSchedule(svc, l)
	res.RegisterExecutor(local)
	return res
}

func InitLocalFuncExeutor(svc service.RankingService) *job.LocalFuncExecutor {
	res := job.NewLocalFuncExecutor()
	// TODO: 注意，这里是内部调用，需要在数据库里插入一条记录，是有 ranking job 负责插入数据库的
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.Job) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		return svc.TopN(ctx)
	})
	return res
}

// CronJob 的改造就是从前端录入，生成一条数据库里的定时任务；
