package main

import (
	"context"

	"google.golang.org/grpc/credentials"

	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
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
	// 自定义认证
	auth := &Auth{
		Appkey:    "testname",
		AppSecret: "testpass",
	}

	cc, err := grpc.Dial(port, grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		log.Panicf("grpc dial error:%v", err)
	}
	defer cc.Close()
	client := pb.NewSearchServiceClient(cc)
	// add Deadline
	ctx, cancle := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(time.Second*5)))
	defer cancle()

	resp, err := client.Search(ctx, &pb.SearchRequest{Request: "grpc"})
	if err != nil {
		// 添加一下状态判断
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}

		log.Printf("search error:%v", err)
	}
	log.Println(resp)
}
