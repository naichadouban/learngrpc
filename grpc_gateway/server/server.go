package server

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/naichadouban/learngrpc/grpc_gateway/proto"
	"github.com/naichadouban/learngrpc/grpc_gateway/util"
	"google.golang.org/grpc"
)

var (
	ServerPort  string
	CertName    string
	CertPemPath string
	CertKeyPath string
	EndPoint    string
)

func Serve() (err error) {
	EndPoint = ":" + ServerPort
	lis, err := net.Listen("tcp", EndPoint)
	if err != nil {
		log.Panicf("tcp listen error:%v", err)
	}
	tlsConfig, err := util.GetTlsConfig(CertPemPath, CertKeyPath)
	if err != nil {
		return err
	}
	srv := createInternalServer(lis, tlsConfig)
	log.Println("grpc and https listen on:", ServerPort)
	if err = srv.Serve(tls.NewListener(lis, tlsConfig)); err != nil {
		log.Printf("ListenAndServe: %v\n", err)
	}
	return err
}
func createInternalServer(conn net.Listener, tlsConfig *tls.Config) *http.Server {
	var opts []grpc.ServerOption
	// grpc server
	creds, err := credentials.NewServerTLSFromFile(CertPemPath, CertKeyPath)
	if err != nil {
		log.Panicf("Failed to create server tls credentials, %v\n", err)
	}
	opts = append(opts, grpc.Creds(creds))
	grpcServer := grpc.NewServer(opts...)
	// register grpc pb
	pb.RegisterHelloWorldServer(grpcServer, NewHelloService())
	// gateway server
	ctx := context.Background()
	dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertName)
	if err != nil {
		log.Panicf("Failed to create client TLS credentials %v\n", err)
	}
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := runtime.NewServeMux()
	// register grpc-gateway pb
	if err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
		log.Panicf("Failed to register gw server:%v\n", err)
	}

	// http 服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	return &http.Server{
		Addr:      EndPoint,
		Handler:   util.GrpcHandlerFunc(grpcServer, mux),
		TLSConfig: tlsConfig,
	}
}
