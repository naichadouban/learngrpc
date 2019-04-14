package server

import (
	"context"

	pb "github.com/naichadouban/learngrpc/grpc_gateway/proto"
)

type HelloService struct{}

func (h *HelloService) SayHelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{
		Message: "test",
	}, nil
}
func NewHelloService() *HelloService {
	return &HelloService{}
}
