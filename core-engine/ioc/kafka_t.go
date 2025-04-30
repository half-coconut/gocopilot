package ioc

import (
	"github.com/spf13/viper"
	"testing"
)

func TestPath(t *testing.T) {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	filePath := viper.GetString("redis.addr")
	t.Log("redis 路径：", filePath)
}

func TestInitKafka(t *testing.T) {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	var cfg Config

	err = viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	t.Log("kafka 配置文件：", cfg.Addrs)

}
