package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/tomy0000000/grpc-reddit/reddit/reddit"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost", "the address to connect to")
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedRedditServer
}

func (s *server) CreatePost(ctx context.Context, in *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	log.Printf("Received: %v", in.GetPost())
	return &pb.CreatePostResponse{Post: in.GetPost()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *addr, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRedditServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
