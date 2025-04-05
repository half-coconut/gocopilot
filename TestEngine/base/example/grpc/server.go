package grpc

import "context"

type Server struct {
	// 加了新的方法，组合就不会报错
	UnimplementedUserServiceServer
}

var _ UserServiceServer = &Server{}

func (s *Server) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		User: &User{
			Id:   123,
			Name: "kitty",
		},
	}, nil
}
