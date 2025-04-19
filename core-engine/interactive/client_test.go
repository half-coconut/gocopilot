package main

import (
	intrv1 "TestCopilot/TestEngine/api/proto/gen/intr/v1"
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestGRPCClient(t *testing.T) {
	cc, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := intrv1.NewInteractiveServiceClient(cc)
	resp, err := client.Get(context.Background(), &intrv1.GetRequest{
		Biz:   "test",
		BizId: 1,
		Uid:   1,
	})
	require.NoError(t, err)
	t.Log(resp.Intr)
}
