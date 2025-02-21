package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	addrs      = []string{"localhost:9094"}
	test_topic = "test_topic"
)

func TestSyncProducer(t *testing.T) {
	// 同步发送
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: test_topic,
		// 消息数据本体
		// json 序列化，转 json, 或 protobuf
		Value: sarama.StringEncoder("hello, this is a message A."),
		Headers: []sarama.RecordHeader{{
			Key:   []byte("trace_id"),
			Value: []byte("123456"),
		}},
		// 之作用于发送过程
		Metadata: "this is metadata",
	})
	assert.NoError(t, err)
}

func TestSyncProducer_partitioner(t *testing.T) {
	// 指定分区发送

	//client,err := sarama.NewClient(addrs,cfg)
	//assert.NoError(t, err)
	//producer,err := sarama.NewSyncProducerFromClient(client)
	//assert.NoError(t, err)

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	// 根据key的hash值 筛选一个
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: test_topic,
		// 消息数据本体
		// json 序列化，转 json, 或 protobuf
		Value: sarama.StringEncoder("hello, this is a message B."),
		Headers: []sarama.RecordHeader{{
			Key:   []byte("trace_id"),
			Value: []byte("123456"),
		}},
		// 之作用于发送过程
		Metadata: "this is metadata",
	})
	assert.NoError(t, err)
}

func Test_Async_producer_acks(t *testing.T) {
	// 异步发送
	cfg := sarama.NewConfig()
	// 关心发送成功和不成功的
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	// 指定 acks
	// 1. 发送一次，不需要服务端的确认
	//cfg.Producer.RequiredAcks = sarama.NoResponse // 0
	// 2. 发送，需要服务端写入主分区
	//cfg.Producer.RequiredAcks = sarama.WaitForLocal // 1
	// 3. 发送，需要服务端同步到所有
	cfg.Producer.RequiredAcks = sarama.WaitForAll // -1

	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	// require 如果有错，会panic
	require.NoError(t, err)

	msg := &sarama.ProducerMessage{
		Topic: test_topic,
		Value: sarama.StringEncoder("hello, this is a message C."),
		Headers: []sarama.RecordHeader{{
			Key:   []byte("trace_id"),
			Value: []byte("123456"),
		}},
		Metadata: "this is metadata",
	}

	// 单向channel 发送通道
	msgCh := producer.Input()
	select {
	case msgCh <- msg:
		//default:
	}

	// 单向channel 接收通道
	errCh := producer.Errors()
	succCh := producer.Successes()

	// select 如果同时满足，就会随机执行一条
	// 如果两种情况都没发生，就会阻塞在 select 这里
	select {
	case err := <-errCh:
		t.Log("发送不成功", err.Err, err.Msg.Value)
	case msg := <-succCh:
		t.Log("发送成功", msg.Value)
	}
}

type JSONEncoder struct {
	Data any
}
