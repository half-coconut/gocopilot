package grpcx

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type Server struct {
	*grpc.Server
	Port     int
	EtcdAddr []string
	Name     string
	L        logger.LoggerV1
	kaCancel func()
	em       endpoints.Manager
	client   *etcdv3.Client
	key      string
}

// 使用 etcd 改造注册方式
// []string{"localhost:12379"}

func (s *Server) Serve() error {
	//l, err := net.Listen("tcp", ":8090")
	l, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		panic(err)
	}
	err = s.register()
	if err != nil {
		return err
	}
	return s.Server.Serve(l)
}

func (s *Server) register() error {
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: s.EtcdAddr,
	})
	if err != nil {
		return err
	}
	s.client = client
	em, err := endpoints.NewManager(client, "service/"+s.Name)
	if err != nil {
		return err
	}

	//addr := "127.0.0.1:8090"
	addr := GetOutBoundIP() + ":" + strconv.Itoa(s.Port)
	key := "service/" + s.Name + "/" + addr
	s.key = key

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 租期可以配置
	var ttl int64 = 30
	leaseResp, err := client.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// 这里 context 使用的是 WithCancel
	kaCtx, kaCancel := context.WithCancel(context.Background())
	s.kaCancel = kaCancel
	ch, err := client.KeepAlive(kaCtx, leaseResp.ID)
	if err != nil {
		return err
	}
	go func() {
		// 续约

		for kaResp := range ch {
			s.L.Debug(kaResp.String())
		}
	}()
	return nil
}

func (s *Server) Close() error {
	if s.kaCancel != nil {
		s.kaCancel()
	}

	if s.em != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.em.DeleteEndpoint(ctx, s.key)
		if err != nil {
			return err
		}
	}
	// 如果是依赖注入，就不要关闭
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			return err
		}
	}
	s.GracefulStop()
	return nil
}
