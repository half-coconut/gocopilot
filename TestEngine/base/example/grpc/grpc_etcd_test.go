package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

type EtcdTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (s *EtcdTestSuite) SetupSuite() {
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(s.T(), err)
	s.client = client
}

func (s *EtcdTestSuite) TestClient() {
	bd, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(bd),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := NewUserServiceClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetById(ctx, &GetByIdReq{
		Id: 1,
	})
	require.NoError(s.T(), err)
	s.T().Log(resp.User)
	//time.Sleep(time.Minute)
}

func (s *EtcdTestSuite) TestServer() {
	l, err := net.Listen("tcp", ":8090")
	require.NoError(s.T(), err)

	em, err := endpoints.NewManager(s.client, "service/user")
	require.NoError(s.T(), err)

	//addr := "127.0.0.1:8090"
	addr := GetOutBoundIP() + ":8090"
	key := "service/user/" + addr

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	var ttl int64 = 30
	leaseResp, err := s.client.Grant(ctx, ttl)
	require.NoError(s.T(), err)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
		// 在这里添加权重
		Metadata: map[string]any{
			"weight": 100,
		},
	}, etcdv3.WithLease(leaseResp.ID))
	require.NoError(s.T(), err)

	// 这里 context 使用的是 WithCancel
	kaCtx, kaCancel := context.WithCancel(context.Background())
	go func() {
		// 续约
		ch, err1 := s.client.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(s.T(), err1)
		for kaResp := range ch {
			s.T().Log(kaResp.String(), time.Now().String())
		}
	}()

	// 注册信息有变动
	go func() {
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)

			err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
				Addr: addr,
				//Metadata: now.String(),
				// 注意：更新，也需要在这里添加权重
				Metadata: map[string]any{
					"weight": 200,
					"time":   now.String(),
				},
				// 注意：更新注册信息时，需要把 lease ID 带上
			}, etcdv3.WithLease(leaseResp.ID))
			if err != nil {
				s.T().Log(err)
			}
			cancel()
		}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{})
	err = server.Serve(l)
	s.T().Log(err)
	// 退出
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	kaCancel()
	err = em.DeleteEndpoint(ctx, key)
	s.client.Close()
	server.GracefulStop()
}

func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}
