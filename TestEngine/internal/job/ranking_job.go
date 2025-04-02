package job

import (
	"TestCopilot/TestEngine/internal/service"
	"context"
	"time"
)

// 对于 Job, GRPC, Web 都是去调用 Service 层的

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func NewRankingJob(svc service.RankingService, timeout time.Duration) *RankingJob {
	// 根据你的数据量，如果要是 7 天内
	return &RankingJob{svc: svc, timeout: timeout}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}
