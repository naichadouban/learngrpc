package main

import (
	"context"
	"log"

	pb "github.com/naichadouban/learngrpc/grpc_gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("../certs/server.pem", "server")
	if err != nil {
		log.Panicf("new client creds error:%v\n", err)
	}
	cc, err := grpc.Dial(":8010", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Panicf("create client connection error:%v\n ", err)
	}
	defer cc.Close()
	c := pb.NewHelloWorldClient(cc)
	context := context.Background()
	body := &pb.HelloWorldRequest{
		Referer: "Grpc",
	}
	res, err := c.SayHelloWorld(context, body)
	if err != nil {
		log.Panicf("grpc client request error:%v\n", err)
	}
	log.Println(res.Message)
}
