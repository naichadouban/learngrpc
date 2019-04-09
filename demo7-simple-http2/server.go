package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/naichadouban/learngrpc/demo7-simple-http/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	certFile := "./conf/server/server.crt"
	keyFile := "./conf/server/server.key"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Panicln(err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./conf/ca.crt")
	if err != nil {
		log.Println(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Panicln("certPool.AppendCertsFromPEM err")
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		// ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientAuth: tls.NoClientCert,
		ClientCAs:  certPool,
	})

	mux := getHTTPServerMux()

	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(server, &SearchService{})

	http.ListenAndServeTLS(port, certFile, keyFile, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("accept")
		fmt.Println(r.Header.Get("Content-Type"))
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
