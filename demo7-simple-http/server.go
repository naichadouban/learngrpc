package main

import (
	"context"

	pb "github.com/naichadouban/learngrpc/demo1/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
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
	mux := getHTTPServerMux()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("listen error:%v", err)
	}
	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})

	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
		return
	}))

}

func getHTTPServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("eddycjy: go-grpc-example"))
	})
	return mux
}
