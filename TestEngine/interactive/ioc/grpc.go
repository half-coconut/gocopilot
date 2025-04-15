package ioc

import (
	grpc2 "TestCopilot/TestEngine/interactive/grpc"
	"TestCopilot/TestEngine/pkg/grpcx"
	"TestCopilot/TestEngine/pkg/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// InitGRPCxServer 类似 web 里的加 handle 一样
func InitGRPCxServer(l logger.LoggerV1, intrServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Port     int      `yaml:"port"`
		EtcdAddr []string `yaml:"etcdAddr"`
	}

	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	// 完成某个微服务的注册，读取配置文件的地址
	server := grpc.NewServer()
	intrServer.Register(server)

	return &grpcx.Server{
		Server:   server,
		Port:     cfg.Port,
		EtcdAddr: cfg.EtcdAddr,
		Name:     "interactive",
		L:        l,
	}
}
