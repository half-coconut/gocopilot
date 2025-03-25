package note

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProducerReadEvent(ctx context.Context, evt ReadEvent) error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

// ProducerReadEvent 如果有复杂的重试逻辑，优先用装饰器，重试逻辑剥离出去，否则就放在这里
func (k KafkaProducer) ProducerReadEvent(ctx context.Context, evt ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_note",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func NewKafkaProducer(pc sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: pc,
	}
}

type ReadEvent struct {
	Uid int64
	Nid int64
}
