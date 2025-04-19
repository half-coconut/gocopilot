package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/half-coconut/gocopilot/core-engine/internal/job"
	"github.com/half-coconut/gocopilot/core-engine/internal/service"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	cronv3 "github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService, client *rlock.Client, l logger.LoggerV1) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30, client, l)

}

func InitJobs(l logger.LoggerV1, rankingJob *job.RankingJob) *cronv3.Cron {
	res := cronv3.New(cronv3.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	// 每 3 分钟执行一次
	_, err := res.AddJob("0 */3 * * * ?", cbd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return res
}
