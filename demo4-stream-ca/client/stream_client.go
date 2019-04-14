package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	pb "github.com/naichadouban/learngrpc/demo4-stream-ca/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
)

const (
	port = ":8011"
)

func main() {
	cert, err := tls.LoadX509KeyPair("../conf/client/client.crt", "../conf/client/client.key")
	if err != nil {
		log.Panicln(err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../conf/ca.crt")
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
		log.Println(err)
	}
	defer cc.Close()
	client := pb.NewStreamServiceClient(cc)
	err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 2018}})
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}

	err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record", Value: 2018}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}

	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 2018}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}

}
func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		log.Println(err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp%v", resp)
	}
	return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		log.Println(err)
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp:%v", resp)
	return nil

}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("%v", resp)
	}
	stream.CloseSend()
	return nil

}
