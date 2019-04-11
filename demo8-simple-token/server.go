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
	"net/http"
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
	if err := s.auth.Check(ctx); err != nil {
		return nil, err
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

	// mux := getHTTPServerMux()

	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}
	err = server.Serve(lis)
	// err = http.ListenAndServeTLS(port, certFile, keyFile, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("accept")
	// 	fmt.Println(r.Header.Get("Content-Type"))
	// 	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
	// 		server.ServeHTTP(w, r)
	// 	} else {
	// 		mux.ServeHTTP(w, r)
	// 	}
	// 	return
	// }))
	log.Println(err)

}

func getHTTPServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("eddycjy: go-grpc-example"))
	})
	return mux
}
