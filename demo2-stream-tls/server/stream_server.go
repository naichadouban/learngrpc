package main

import (
	pb "github.com/naichadouban/learngrpc/demo2-stream/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net"
)

const (
	port = ":8011"
)

type StreamService struct {
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
		}
		if err != nil {
			log.Println(err)
		}
		log.Printf("stream recv:%v", r)
	}
	return nil
}
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "grpc stream clinet:route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("stream recv %v", r)
	}
	return nil
}

func main() {
	c, err := credentials.NewServerTLSFromFile("../conf/server.pem", "../conf/server.key")
	if err != nil {
		log.Panicln(err)
	}
	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}
	server.Serve(lis)
}
