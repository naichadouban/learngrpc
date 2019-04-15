package main

import (
	"context"
	"log"

	pb "github.com/naichadouban/learngrpc/demo6-interceptor/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8010"
)

func main() {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc dial error:%v", err)
	}
	defer conn.Close()
	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{Request: "grpc"})
	if err != nil {
		log.Printf("search error:%v", err)
	}
	log.Println(resp)
}
