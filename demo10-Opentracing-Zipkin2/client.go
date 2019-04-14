package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/openzipkin-contrib/zipkin-go-opentracing"

	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"google.golang.org/grpc"

	"log"
)

const (
	port                      = ":8010"
	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://192.168.1.74:9411/api/v2/spans"
	ZIPKIN_RECORDER_HOST_PORT = "192.168.1.74:9000"
)

func main() {
	// opentrace
	collector, err := zipkintracer.NewHTTPCollector(ZIPKIN_HTTP_ENDPOINT)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err: %v", err)
	}
	recorder := zipkintracer.NewRecorder(collector, true, ZIPKIN_RECORDER_HOST_PORT, SERVICE_NAME)
	tracer, err := zipkintracer.NewTracer(
		recorder, zipkintracer.ClientServerSameSpan(true),
	)
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}
	// ca tls
	//cert, err := tls.LoadX509KeyPair("./conf/client/client.crt", "./conf/client/client.key")
	//if err != nil {
	//	log.Panicln(err)
	//}
	//certPool := x509.NewCertPool()
	//ca, err := ioutil.ReadFile("./conf/ca.crt")
	//if err != nil {
	//	log.Panicln(err)
	//}
	//if ok := certPool.AppendCertsFromPEM(ca); !ok {
	//	log.Println("certPool.AppendCertsFromPEM err")
	//}
	//c := credentials.NewTLS(&tls.Config{
	//	Certificates: []tls.Certificate{cert},
	//	ServerName:   "server",
	//	RootCAs:      certPool,
	//})
	// connect
	cc, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithUnaryInterceptor(
		otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()),
	))
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
