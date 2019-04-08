package main

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/naichadouban/learngrpc/demo1/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime/debug"
)

type SearchService struct {
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + "server"}, nil
}

const (
	port = ":8010"
)

func LogginInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}
func main() {
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			LogginInterceptor,
			RecoveryInterceptor,
		),
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("listen error:%v", err)
	}
	server := grpc.NewServer(opts...)

	pb.RegisterSearchServiceServer(server, &SearchService{})

	server.Serve(lis)
}
