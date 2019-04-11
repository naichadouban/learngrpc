package main

import (
	"context"
	"log"

	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

const (
	port = ":8010"
)

type Auth struct {
	appKey    string
	appSecret string
}

func (a *Auth) GetAppKey() string {
	return "testname"
}

func (a *Auth) GetAppSecret() string {
	return "testpass"
}
func (a *Auth) Check(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "自定义认证 Token 失败")
	}

	var (
		appKey    string
		appSecret string
	)
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}

	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return status.Errorf(codes.Unauthenticated, "自定义认证 Token 无效")
	}

	return nil
}

type SearchService struct {
	auth *Auth
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	// 这是自定义的认证
	if err := s.auth.Check(ctx); err != nil {
		return nil, err
	}
	for i := 0; i < 3; i++ {
		// 而在 Server 端，由于 Client 已经设置了截止时间。Server 势必要去检测它
		// 否则如果 Client 已经结束掉了，Server 还傻傻的在那执行，这对资源是一种极大的浪费
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "SearchService.Search canceled")
		}

		time.Sleep(1 * time.Second)
		log.Printf("time:%d\n", i)
	}
	return &pb.SearchResponse{Response: r.GetRequest() + "server"}, nil
}

func main() {
	certFile := "./conf/server/server.crt"
	keyFile := "./conf/server/server.key"
	c, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Panicln(err)
	}

	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}
	err = server.Serve(lis)
	log.Println(err)

}
