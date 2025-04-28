package grpc

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/pkg/ginx/balancer/wrr"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	_ "google.golang.org/grpc/balancer/weightedroundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

const name = "custom_wrr"

func init() {
	//NewBalancerBuilder 把 pickerbuilder 转化为注册一个 balancer.Builder
	balancer.Register(base.NewBalancerBuilder(name, &wrr.PickerBuilder{}, base.Config{}))
}

type BalanceTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

//func (s *BalanceTestSuite) TestPickFirst(t *testing.T) {
//
//}

func (s *BalanceTestSuite) SetupSuite() {
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(s.T(), err)
	s.client = client
}

func TestBalanceTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}

func (s *BalanceTestSuite) TestRoundRobinClient() {
	bd, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	svcCfg := `
{
    "loadBalancingConfig": [
        {
            "round_robin": {}
        }
    ]
}
`
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(bd),
		// 在这里使用的负载均衡器
		grpc.WithDefaultServiceConfig(svcCfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := NewUserServiceClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.GetById(ctx, &GetByIdReq{
		Id: 1,
	})
	cancel()
	require.NoError(s.T(), err)
	s.T().Log(resp.User)
}

func (s *BalanceTestSuite) TestWeightedRoundRobinClient() {
	bd, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	svcCfg := `
{
    "loadBalancingConfig": [
        {
            "weighted_round_robin": {}
        }
    ]
}
`
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(bd),
		grpc.WithDefaultServiceConfig(svcCfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := NewUserServiceClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.GetById(ctx, &GetByIdReq{
		Id: 1,
	})
	cancel()
	require.NoError(s.T(), err)
	s.T().Log(resp.User)
}

func (s *BalanceTestSuite) TestCustomWRRClient() {
	bd, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	svcCfg := `
{
    "loadBalancingConfig": [
        {
            "custom_wrr": {}
        }
    ]
}
`
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(bd),
		grpc.WithDefaultServiceConfig(svcCfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := NewUserServiceClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	resp, err := client.GetById(ctx, &GetByIdReq{
		Id: 1,
	})
	cancel()
	require.NoError(s.T(), err)
	s.T().Log(resp.User)
}
