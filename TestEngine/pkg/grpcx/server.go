package grpcx

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	Addr string
}

func (s *Server) Serve(addr string) error {
	//l, err := net.Listen("tcp", ":8090")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	return s.Server.Serve(l)
}
