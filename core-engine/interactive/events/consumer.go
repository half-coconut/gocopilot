package note

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/half-coconut/gocopilot/core-engine/interactive/repository"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"github.com/half-coconut/gocopilot/core-engine/pkg/saramax"
	"time"
)

var _ saramax.Consumer = &InteractiveReadEventConsumer{}

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveReadEventConsumer(l logger.LoggerV1, repo repository.InteractiveRepository, client sarama.Client) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (k *InteractiveReadEventConsumer) Start() error {
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
func (k *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return k.repo.IncrReadCnt(ctx, "note", t.Nid)
}
