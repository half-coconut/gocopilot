package sarama

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	log2 "log"
	"testing"
	"time"
)

// 需要知道何时结束，使用context，控制消费者退出：
// context.WithTimeout
// context.WithCancel

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_group", cfg)
	assert.NoError(t, err)

	err = consumer.Consume(context.Background(), []string{test_topic}, testConsumerGroupHandler{})
	assert.NoError(t, err)
}

func TestConsumer_With_Timeout(t *testing.T) {
	// 5秒超时后，结束
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_group", cfg)
	assert.NoError(t, err)

	begin := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = consumer.Consume(ctx, []string{test_topic}, testConsumerGroupHandler{})
	assert.NoError(t, err)
	t.Log(time.Since(begin).String())
}

func TestConsumer_With_Cancel(t *testing.T) {
	// 30秒后，结束
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_group", cfg)
	assert.NoError(t, err)

	begin := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*30, func() {
		cancel()
	})
	err = consumer.Consume(ctx, []string{test_topic}, testConsumerGroupHandler{})
	assert.NoError(t, err)
	t.Log(time.Since(begin).String())
}

type testConsumerGroupHandler struct {
}

func (c testConsumerGroupHandler) SetupV1(session sarama.ConsumerGroupSession) error {
	log2.Println("Setup")
	return nil
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// 在 setup 里指定偏移量消费
	// topic => 偏移量
	partitions := session.Claims()[test_topic]

	// 一般保留 3-7天
	for _, part := range partitions {
		// 从最晚的开始去消费
		session.ResetOffset(test_topic, part,
			sarama.OffsetOldest, "")
		// 或者是 offset 直接填写，比如 123
		//session.ResetOffset("test_topic", part,
		//	123, "")
		//session.ResetOffset("test_topic", part,
		//	sarama.OffsetNewest, "")
	}

	return nil
}

func (c testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log2.Println("Cleanup")
	return nil
}

type MyBizMsg struct {
}

func (c testConsumerGroupHandler) ConsumeClaimV3(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgCh := claim.Messages()
	for msg := range msgCh {
		var bisMsg MyBizMsg
		err := json.Unmarshal(msg.Value, &bisMsg)
		log2.Println(string(msg.Value))
		if err != nil {
			continue
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c testConsumerGroupHandler) ConsumeClaimV2(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgCh := claim.Messages()
	for msg := range msgCh {
		log2.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c testConsumerGroupHandler) ConsumeClaim(
	// 代表的是你和Kafka 的会话（从建立连接到连接彻底断掉的那一段时间）
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgCh := claim.Messages()

	const batchSize = 10
	for {
		// 这里用于控制凑够一批的时间
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		var (
			eg   errgroup.Group
			last *sarama.ConsumerMessage
		)
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 超时了后，退出循环
				done = true
			case msg, ok := <-msgCh:
				if !ok {
					cancel()
					// 消费者 msgCh 被关闭了
					return nil
				}
				last = msg
				eg.Go(func() error {
					// 在这里消费，重试
					time.Sleep(time.Second)
					log2.Println(string(msg.Value))
					return nil
				})
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			// 记录日志
			// 这里是一整批重试
			continue
		}
		// 做实验，是否是最后一个msg 就生效
		if last != nil {
			session.MarkMessage(last, "")
		}
	}
}

// 返回只读的 channel
func ChannelV1() <-chan struct{} {
	panic("implement me")
}

// 返回只写的 channel
func ChannelV2() chan<- struct{} {
	panic("implement me")
}

// 返回可读可写的 channel
func ChannelV3() chan struct{} {
	panic("implement me")
}
