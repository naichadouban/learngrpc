package main

import (
	"context"
	pb "github.com/naichadouban/learngrpc/demo1/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type SearchService struct {
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + "server"}, nil
}

const (
	port = ":8010"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("listen error:%v", err)
	}
	server := grpc.NewServer()

	pb.RegisterSearchServiceServer(server, &SearchService{})

	server.Serve(lis)
}
