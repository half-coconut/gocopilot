package redisx

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

func NewPrometheusHook(opt prometheus.SummaryOpts) *PrometheusHook {
	vector := prometheus.NewSummaryVec(opt,
		// 是否命中缓存
		[]string{"cmd", "key_exist"})
	prometheus.MustRegister(vector)
	return &PrometheusHook{
		vector: vector,
	}
}

func (p *PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 相当于什么也没干
		return next(ctx, network, addr)
	}
}

func (p *PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	// 在这里监控
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 在 redis 执行之前
		startTime := time.Now()
		var err error
		defer func() {
			duration := time.Since(startTime).Milliseconds()
			//biz := ctx.Value("biz")
			keyExist := err == redis.Nil
			p.vector.WithLabelValues(
				cmd.Name(),

				strconv.FormatBool(keyExist)).Observe(float64(duration))
		}()
		err = next(ctx, cmd)
		// 在 redis 执行之后
		return err
	}

}

func (p *PrometheusHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	//TODO implement me
	panic("implement me")
}
