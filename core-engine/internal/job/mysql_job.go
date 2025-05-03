package job

//
//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
//	"github.com/half-coconut/gocopilot/core-engine/internal/service"
//	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
//	"golang.org/x/sync/semaphore"
//	"net/http"
//	"time"
//)
//
//type Executor interface {
//	Name() string
//	Exec(ctx context.Context, j domain.Job) error
//}
//
//// HttpExecutor 执行
//type HttpExecutor struct {
//}
//
//func (h HttpExecutor) Name() string {
//	return "http"
//}
//
//func (h HttpExecutor) Exec(ctx context.Context, j domain.Job) error {
//	// 这部分实现是说：任务调度节点，到时间之后，就通过 http 去找你的服务节点的定时任务，然后去执行
//	type Config struct {
//		Endpoint string
//		Method   string
//	}
//	var cfg Config
//	err := json.Unmarshal([]byte(j.Cfg), &cfg)
//	if err != nil {
//		return err
//	}
//	req, err := http.NewRequest(cfg.Method, cfg.Endpoint, nil)
//	if err != nil {
//		return err
//	}
//	resp, _ := http.DefaultClient.Do(req)
//	if resp.StatusCode != http.StatusOK {
//		return errors.New("执行失败")
//	}
//	return nil
//}
//
//// LocalFuncExecutor 本地方法调用
//type LocalFuncExecutor struct {
//	funcs map[string]func(ctx context.Context, j domain.Job) error
//}
//
//func NewLocalFuncExecutor() *LocalFuncExecutor {
//	return &LocalFuncExecutor{
//		funcs: make(map[string]func(ctx context.Context, j domain.Job) error)}
//}
//
//func (l *LocalFuncExecutor) Name() string {
//	return "local"
//}
//func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
//	l.funcs[name] = fn
//}
//
//func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
//	fn, ok := l.funcs[j.Name]
//	if !ok {
//		return fmt.Errorf("未知任务，你是否注册了？%s", j.Name)
//	}
//	return fn(ctx, j)
//}
//
//// Schedule 调度器, mysql-job
//type Schedule struct {
//	svc     service.JobService
//	l       logger.LoggerV1
//	execs   map[string]Executor
//	limiter *semaphore.Weighted
//}
//
//func NewSchedule(svc service.JobService, l logger.LoggerV1) *Schedule {
//	return &Schedule{
//		svc: svc,
//		l:   l,
//		// 使用信号量，最多开出 200 个 goroutine
//		limiter: semaphore.NewWeighted(200),
//		execs:   make(map[string]Executor)}
//}
//
//func (s *Schedule) RegisterExecutor(exec Executor) {
//	s.execs[exec.Name()] = exec
//}
//
//func (s *Schedule) Schedule(ctx context.Context) error {
//	// 调度器，是一个 for 循环，一直开启的服务
//	for {
//
//		if ctx.Err() != nil {
//			// 退出调度循环
//			return ctx.Err()
//		}
//		err := s.limiter.Acquire(ctx, 1)
//		if err != nil {
//			return err
//		}
//		// 一次调度数据库查询时间
//		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
//		defer cancel()
//
//		j, err := s.svc.Preempt(dbCtx)
//		if err != nil {
//			s.l.Error("抢占任务失败", logger.Error(err))
//		}
//
//		exec, ok := s.execs[j.Executor]
//		if !ok {
//			s.l.Error("未找到对应的执行器", logger.String("executor", j.Executor))
//			continue
//		}
//		// 抢完了就执行，先Preempt
//		// 执行完毕
//		go func() {
//			defer func() {
//				s.limiter.Release(1)
//				er1 := j.CancelFunc()
//				if er1 != nil {
//					s.l.Error("释放任务失败",
//						logger.Int64("jid", j.Id),
//						logger.Error(er1))
//				}
//			}()
//
//			er1 := exec.Exec(ctx, j)
//			if er1 != nil {
//				// 可以重试
//				s.l.Error("任务执行失败", logger.Error(er1))
//			}
//			// 考虑下一次调度
//			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
//			defer cancel()
//			er1 = s.svc.ResetNextTime(ctx, j)
//			if er1 != nil {
//				s.l.Error("设置下一次执行时间失败", logger.Error(er1))
//			}
//		}()
//	}
//}
