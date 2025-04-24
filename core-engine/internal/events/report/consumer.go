package report

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"github.com/half-coconut/gocopilot/core-engine/pkg/saramax"
	"time"
)

var _ saramax.Consumer = &DebugLogEventConsumer{}

type DebugLogEventConsumer struct {
	client sarama.Client
	repo   repository.ReportRepository
	l      logger.LoggerV1
}

func NewDebugLogEventConsumer(client sarama.Client, repo repository.ReportRepository, l logger.LoggerV1) *DebugLogEventConsumer {
	return &DebugLogEventConsumer{client: client, repo: repo, l: l}
}

func (d *DebugLogEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("report",
		d.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicDebugLogEvent},
			saramax.NewHandler[domain.DebugLog](d.l, d.Consume))
		if err != nil {
			d.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (d *DebugLogEventConsumer) Consume(msg *sarama.ConsumerMessage, log domain.DebugLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := d.repo.CreateDebugLog(ctx, log)
	d.l.Info(fmt.Sprintf("通过消费者消费，存入数据库: %v", log))
	return err
}
