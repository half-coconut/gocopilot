package ioc

import (
	"github.com/IBM/sarama"
	events "github.com/half-coconut/gocopilot/core-engine/internal/events/report"
	"github.com/half-coconut/gocopilot/core-engine/pkg/saramax"
	"github.com/spf13/viper"
	"log"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true

	//viper.SetConfigName("dev")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("../config")
	//err := viper.ReadInConfig()
	var cfg Config

	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

// NewConsumers 注意：所有的 Consumer 在这里注册一下
func NewConsumers(debugLog *events.DebugLogEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{
		debugLog,
	}
}
