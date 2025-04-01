package http

import (
	"context"
	"fmt"
	rate2 "golang.org/x/time/rate"
	"log"
	"sync"
	"time"
)

func (t *Target) http_load(duration time.Duration, rate float64) *Result {
	// 创建限速器
	limiter := rate2.NewLimiter(rate2.Limit(rate), 1)

	// 创建上下文，用于控制压测持续时间
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// 创建 WaitGroup，用于等待所有请求完成
	var wg sync.WaitGroup
	var res *Result

	// 启动多个 goroutine 发送请求

	for {
		select {
		case <-ctx.Done():
			// 压测时间到，退出循环
			fmt.Println("压测结束")
			return res
		default:
			// 判断是否获取到令牌
			if limiter.Allow() {
				// 启动 goroutine 发送请求

				wg.Add(1)
				go func() {
					defer wg.Done()

					// 发送 HTTP 请求
					res = t.Do()
				}()
			} else {
				// 处理限速，例如记录日志或等待一段时间
				log.Printf("限速，当前的 rate limit 为：%v\n", limiter.Limit())
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

}
