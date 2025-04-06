package ioc

import (
	intrv1 "TestCopilot/TestEngine/api/proto/gen/intr/v1"
	"TestCopilot/TestEngine/interactive/service"
	"TestCopilot/TestEngine/internal/web/client"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func InitIntrGRPCClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string
		Secure    bool
		Threshold int32
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if cfg.Secure {
		// 加载证书
		// 启用 HTTPS
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewInteractiveServiceAdapter(svc)
	res := client.NewGreyScaleInteractiveServiceClient(remote, local)

	// 在这里监听
	viper.OnConfigChange(func(in fsnotify.Event) {
		err = viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {
			// 输入日志
			log.Println("监听配置变更失败", err)
		}
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
