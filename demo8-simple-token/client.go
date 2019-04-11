package main

import (
	"context"

	"google.golang.org/grpc/credentials"

	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"google.golang.org/grpc"

	"log"
)

const (
	port = ":8010"
)

type Auth struct {
	Appkey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.Appkey, "app_secret": a.AppSecret}, nil
}
func (a *Auth) RequireTransportSecurity() bool {
	return true
}

func main() {
	c, err := credentials.NewClientTLSFromFile("./conf/server/server.crt", "server")
	if err != nil {
		log.Println(err)
	}
	auth := &Auth{
		Appkey: "testname",
		//AppSecret: "testpass",
		AppSecret: "testpass",
	}
	cc, err := grpc.Dial(port, grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(auth))
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
