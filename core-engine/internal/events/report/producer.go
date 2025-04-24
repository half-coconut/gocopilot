package report

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/zeromicro/go-zero/core/jsonx"
)

const topicDebugLogEvent = "debug_logs_event"

type DebugLogProducer interface {
	ProducerRecordDebugLogsEvent(ctx context.Context, evt domain.DebugLog) error
}

type KafkaDebugLogProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaDebugLogProducer(producer sarama.SyncProducer) DebugLogProducer {
	return &KafkaDebugLogProducer{producer: producer}
}

func (k *KafkaDebugLogProducer) ProducerRecordDebugLogsEvent(ctx context.Context, evt domain.DebugLog) error {
	data, err := jsonx.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicDebugLogEvent,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
