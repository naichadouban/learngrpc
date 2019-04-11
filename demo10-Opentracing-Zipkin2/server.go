package main

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	"log"

	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"github.com/openzipkin-contrib/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"net"
)

const (
	port                      = ":8010"
	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://192.168.1.74:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "192.168.1.74:9000"
)

type SearchService struct {
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + "server"}, nil
}

func main() {
	// opentrace
	collector, err := zipkintracer.NewHTTPCollector(ZIPKIN_HTTP_ENDPOINT)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err: %v", err)
	}
	recorder := zipkintracer.NewRecorder(collector, true, ZIPKIN_RECORDER_HOST_PORT, SERVICE_NAME)
	tracer, err := zipkintracer.NewTracer(
		recorder, zipkintracer.ClientServerSameSpan(false),
	)
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}
	// ca tls
	//cert, err := tls.LoadX509KeyPair("./conf/server/server.crt", "./conf/server/server.key")
	//if err != nil {
	//	log.Panicln(err)
	//}
	//certPool := x509.NewCertPool()
	//ca, err := ioutil.ReadFile("./conf/ca.crt")
	//if err != nil {
	//	log.Println(err)
	//}
	//if ok := certPool.AppendCertsFromPEM(ca); !ok {
	//	log.Panicln("certPool.AppendCertsFromPEM err")
	//}
	//c := credentials.NewTLS(&tls.Config{
	//	Certificates: []tls.Certificate{cert},
	//	ClientAuth:   tls.RequireAndVerifyClientCert,
	//	ClientCAs:    certPool,
	//})

	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}
	err = server.Serve(lis)
	log.Println(err)

}
