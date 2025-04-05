package interactive

import (
	intrv1 "TestCopilot/TestEngine/api/proto/gen/intr/v1"
	"TestCopilot/TestEngine/interactive/grpc"
	grpc2 "google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	server := grpc2.NewServer()
	intrSvc := &grpc.InteractiveServiceServer{}
	intrv1.RegisterInteractiveServiceServer(server, intrSvc)
	// 监听 8090 端口
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	// 这边会阻塞，类似于 gin.Run
	err = server.Serve(l)
	log.Println(err)
}
