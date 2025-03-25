package note

import (
	"TestCopilot/TestEngine/internal/repository"
	"TestCopilot/TestEngine/pkg/logger"
	"TestCopilot/TestEngine/pkg/saramax"
	"context"
	"github.com/IBM/sarama"
	"time"
)

type Consumer interface {
	Start() error
}

type KafkaConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewKafkaConsumer(l logger.LoggerV1, repo repository.InteractiveRepository, client sarama.Client) *KafkaConsumer {
	return &KafkaConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (k *KafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}

	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_note"},
			saramax.NewHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			k.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume 这个不是幂等
func (k *KafkaConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return k.repo.IncrReadCnt(ctx, "note", t.Nid)
}
