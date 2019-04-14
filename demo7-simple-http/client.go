package main

import (
	"context"
	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"google.golang.org/grpc"
	"log"
)

const (
	port = ":8010"
)

func main() {
	cc, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc dial error:%v", err)
	}
	defer cc.Close()
	client := pb.NewSearchServiceClient(cc)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{Request: "grpc"})
	if err != nil {
		log.Printf("search error:%v", err)
	}
	log.Println(resp)
}
