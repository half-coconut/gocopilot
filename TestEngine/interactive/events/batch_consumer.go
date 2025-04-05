package note

import (
	"TestCopilot/TestEngine/interactive/repository"
	"TestCopilot/TestEngine/pkg/logger"
	"TestCopilot/TestEngine/pkg/saramax"
	"context"
	"github.com/IBM/sarama"
	"time"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveReadEventBatchConsumer(client sarama.Client, repo repository.InteractiveRepository, l logger.LoggerV1) *InteractiveReadEventBatchConsumer {
	return &InteractiveReadEventBatchConsumer{client: client, repo: repo, l: l}
}

func (k *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}

	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_note"},
			saramax.NewBatchHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			k.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume 这个不是幂等
func (k *InteractiveReadEventBatchConsumer) Consume(msgs []*sarama.ConsumerMessage, ts []ReadEvent) error {

	noteIds := make([]int64, 0, len(ts))
	bizs := make([]string, 0, len(ts))
	for _, evt := range ts {
		noteIds = append(noteIds, evt.Nid)
		bizs = append(bizs, "note")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := k.repo.BatchIncrReadCnt(ctx, bizs, noteIds)
	if err != nil {
		k.l.Error("批量增加阅读计数失败",
			logger.Field{Key: "ids", Value: noteIds},
			logger.Error(err))
	}
	return nil
}
