package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	// cc 是连接池的连接池，cc 里面放了很多个连接池
	cc, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials())) // 没有使用 https的时候
	require.NoError(t, err)
	client := NewUserServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetById(ctx, &GetByIdReq{
		Id: 456,
	})
	assert.NoError(t, err)
	t.Log(resp.User)
}
