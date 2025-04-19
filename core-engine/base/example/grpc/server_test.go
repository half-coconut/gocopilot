package grpc

import (
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	// grpc 的 server
	server := grpc.NewServer()
	defer func() {
		// 优雅退出
		server.GracefulStop()
	}()
	// 我们的业务的 server
	userServer := &Server{}
	RegisterUserServiceServer(server, userServer)
	// 创建一个监听器，监听 tcp 协议，8090 端口
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	t.Log(err)

	// 可以使用 http 调用，用于兼容历史老系统
	//server.ServeHTTP()

}
