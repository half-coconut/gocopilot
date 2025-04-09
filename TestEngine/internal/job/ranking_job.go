package job

import (
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"sync"
	"time"
)

// 对于 Job, GRPC, Web 都是去调用 Service 层的

type RankingJob struct {
	svc       service.RankingService
	timeout   time.Duration
	client    *rlock.Client // 使用 redis 分布式锁，在 job 层面加锁
	key       string
	l         logger.LoggerV1
	lock      *rlock.Lock
	localLock *sync.Mutex
}

func NewRankingJob(svc service.RankingService, timeout time.Duration, client *rlock.Client, l logger.LoggerV1) *RankingJob {
	// 根据你的数据量，如果要是 7 天内
	return &RankingJob{
		svc:       svc,
		timeout:   timeout,
		client:    client,
		key:       "rlock:cron_job:ranking",
		l:         l,
		localLock: &sync.Mutex{},
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

// Run 按时间调度
func (r *RankingJob) Run() error {
	r.localLock.Lock()
	defer r.localLock.Unlock()

	if r.lock == nil {
		// 说明没有拿到锁，你试着拿锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// 分布式锁的过期时间，设置一个比较短的时间
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			// 固定间隔重试
			Interval: time.Millisecond * 100,
			// 重试次数
			Max: 0,
		}, time.Second) // 单一 一次调用 redis 的超时时间
		if err != nil {
			// 没有拿到锁
			return nil
		}
		r.lock = lock
		// 需要保证一直拿着这个锁
		go func() {
			// 自动续约机制，AutoRefresh 是阻塞的
			er := lock.AutoRefresh(r.timeout/2, time.Second)
			// 说明退出了续约机制
			if er != nil {
				r.l.Error("续约失败", logger.Error(err))
			}
			r.localLock.Lock()
			r.lock = nil
			r.localLock.Unlock()
		}()
	}

	//defer func() {
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//	defer cancel()
	//	err = lock.Unlock(ctx)
	//	if err != nil {
	//		r.l.Error("释放分布式锁失败, Ranking Job", logger.Error(err))
	//	}
	//}()

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.lock = nil
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}
