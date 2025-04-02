package ioc

import (
	"TestCopilot/TestEngine/internal/job"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/pkg/logger"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30)
}

func InitJobs(l logger.LoggerV1, rankingJob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	// 每 3 分钟执行一次
	_, err := res.AddJob("0 */3 * * * ?", cbd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return res
}
