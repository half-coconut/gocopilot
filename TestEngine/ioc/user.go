package ioc

import (
	"TestCopilot/TestEngine/internal/repository/cache"
	"TestCopilot/TestEngine/pkg/redisx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// InitUserHook 配合 PrometheusHook 使用
func InitUserHook(client *redis.Client) cache.UserCache {
	client.AddHook(redisx.NewPrometheusHook(prometheus.SummaryOpts{
		Namespace: "test_copilot",
		Subsystem: "test_engine",
		Name:      "gin_http",
		Help:      "分业务监控 redis 缓存",
		ConstLabels: map[string]string{
			"biz": "user",
		},
	}))
	panic("先不使用了")
}
