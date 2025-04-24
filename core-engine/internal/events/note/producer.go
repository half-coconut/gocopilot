package note

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type NoteProducer interface {
	ProducerReadEvent(ctx context.Context, evt ReadEvent) error
	ProducerReadEventV1(ctx context.Context, v1 ReadEventV1)
}

type KafkaNoteProducer struct {
	producer sarama.SyncProducer
}

func (k *KafkaNoteProducer) ProducerReadEventV1(ctx context.Context, v1 ReadEventV1) {
	//TODO implement me
	panic("implement me")
}

// ProducerReadEvent 如果有复杂的重试逻辑，优先用装饰器，重试逻辑剥离出去，否则就放在这里
func (k *KafkaNoteProducer) ProducerReadEvent(ctx context.Context, evt ReadEvent) error {
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

func NewKafkaNoteProducer(pc sarama.SyncProducer) NoteProducer {
	return &KafkaNoteProducer{
		producer: pc,
	}
}

type ReadEvent struct {
	Uid int64
	Nid int64
}

type ReadEventV1 struct {
	Uids []int64
	Nids []int64
}
