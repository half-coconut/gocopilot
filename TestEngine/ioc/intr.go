package ioc

import (
	intrv1 "TestCopilot/TestEngine/api/proto/gen/intr/v1"
	"TestCopilot/TestEngine/interactive/service"
	"TestCopilot/TestEngine/internal/web/client"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func InitEtcd() *clientv3.Client {
	var cfg clientv3.Config

	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

// 真正的 gRPC的客户端
func InitIntrGRPCClientV1(client *clientv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr   string
		Name   string
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}

	bd, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}

	opts := []grpc.DialOption{grpc.WithResolvers(bd)}

	if cfg.Secure {

	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial("etcd:///service/"+cfg.Name, opts...)
	if err != nil {
		panic(err)
	}
	return intrv1.NewInteractiveServiceClient(cc)
}

// InitIntrGRPCClient 这是流量控制的客户端
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
