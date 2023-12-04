package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"

	pb "github.com/tomy0000000/grpc-reddit/reddit/reddit"
)

var (
	addr = flag.String("addr", "localhost", "the address to connect to")
	port = flag.Int("port", 50051, "The server port")
)

type gRPCserver struct {
	pb.UnimplementedRedditServer
	sqlClient *SQLClient
}

func (s *gRPCserver) CreatePost(ctx context.Context, in *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	postIn := in.GetPost()
	log.Printf("Received: %v", postIn)

	// Insert the post into the database
	id, err := s.sqlClient.CreatePost(postIn)
	if err != nil {
		log.Fatalf("Error inserting into database: %v", err)
	}

	log.Printf("Post created with id: %d", id)

	// Get the post from the database
	post, err := s.sqlClient.GetPost(id)
	if err != nil {
		log.Fatalf("Error getting post from database: %v", err)
	}
	return &pb.CreatePostResponse{Post: post}, nil
}

func (s *gRPCserver) GetPost(ctx context.Context, in *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	log.Printf("GetPost: %v", in)
	id := in.GetPostID()

	// Get the post from the database
	post, err := s.sqlClient.GetPost(int(id))
	if err != nil {
		log.Fatalf("Error getting post from database: %v", err)
	}
	return &pb.GetPostResponse{Post: post}, nil
}

func main() {
	// Parse the flags
	flag.Parse()

	// Open the database
	s := &gRPCserver{}
	var err error
	s.sqlClient, err = NewSQLClient()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Launch the server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *addr, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	pb.RegisterRedditServer(gs, s)
	log.Printf("server listening at %v", lis.Addr())
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
