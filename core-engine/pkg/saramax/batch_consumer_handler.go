package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

type BatchHandler[T any] struct {
	l  logger.LoggerV1
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	// 用 option 模式来设置 batchSize 和 batchDuration
	batchSize     int
	batchDuration time.Duration
}

func NewBatchHandler[T any](l logger.LoggerV1, fn func(msgs []*sarama.ConsumerMessage, ts []T) error, opts ...Option[T]) *BatchHandler[T] {
	bh := &BatchHandler[T]{l: l, fn: fn, batchSize: 10, batchDuration: time.Second}

	for _, opt := range opts {
		opt(bh) // 应用每个选项
	}

	return bh
}

type Option[T any] func(*BatchHandler[T])

func WithBatchSize[T any](size int) Option[T] {
	return func(bh *BatchHandler[T]) {
		bh.batchSize = size
	}
}

func WithBatchDuration[T any](duration time.Duration) Option[T] {
	return func(bh *BatchHandler[T]) {
		bh.batchDuration = duration
	}
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		done := false
		msgs := make([]*sarama.ConsumerMessage, 0, b.batchSize)
		ts := make([]T, 0, b.batchSize)
		for i := 0; i < b.batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败", logger.Error(err),
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset))
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := b.fn(msgs, ts)
		if err != nil {
			b.l.Error("调用业务批量接口失败", logger.Error(err))
			// 还要继续往前消费
		}
		// 不使用 last，确保数据都对
		for _, msg := range msgs {
			session.MarkMessage(msg, "")
		}
	}
}
