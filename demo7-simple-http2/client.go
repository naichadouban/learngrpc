package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"

	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"google.golang.org/grpc"

	"log"
)

const (
	port = ":8010"
)

func main() {
	certFile := "./conf/client/client.crt"
	keyFile := "./conf/client/client.key"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Panicln(err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./conf/ca.crt")
	if err != nil {
		log.Panicln(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Println("certPool.AppendCertsFromPEM err")
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "server",
		RootCAs:      certPool,
	})
	cc, err := grpc.Dial(port, grpc.WithTransportCredentials(c))
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
