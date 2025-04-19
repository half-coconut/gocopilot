package note

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"github.com/half-coconut/gocopilot/core-engine/pkg/saramax"
)

type HistoryReadEventConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
}

func NewHistoryReadEventConsumer(l logger.LoggerV1, client sarama.Client) *HistoryReadEventConsumer {
	return &HistoryReadEventConsumer{
		client: client,
		l:      l,
	}
}

func (k *HistoryReadEventConsumer) Start() error {
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
func (k *HistoryReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//return k.repo.AddRecord(ctx, t.Nid, t.Nid)
	panic("implement me")
}
