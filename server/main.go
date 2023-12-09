package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
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
	log.Print(color.YellowString("[CreatePost] Received: %v", in))

	// Insert the post into the database
	id, err := s.sqlClient.CreatePost(in.GetPost())
	if err != nil {
		log.Fatal(color.RedString("[CreatePost] DB error: %v", err))
	}

	// Get the post from the database
	post, err := s.sqlClient.GetPost(id)
	if err != nil {
		log.Fatal(color.RedString("[CreatePost] DB error: %v", err))
	}

	response := &pb.CreatePostResponse{Post: post}
	log.Print(color.GreenString("[CreatePost] Reponse: %v", response))
	return response, nil
}

func (s *gRPCserver) VotePost(ctx context.Context, in *pb.VotePostRequest) (*pb.VotePostResponse, error) {
	log.Print(color.YellowString("[VotePost] Received: %v", in))

	// Increment/Decrement the score of the post
	newScore, err := s.sqlClient.VotePost(int(in.GetPostID()), in.GetUpvote())
	if err != nil {
		log.Fatal(color.RedString("[VotePost] DB error: %v", err))
	}

	response := &pb.VotePostResponse{Score: int32(newScore)}
	log.Print(color.GreenString("[VotePost] Reponse: %v", response))
	return response, nil
}

func (s *gRPCserver) GetPost(ctx context.Context, in *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	log.Print(color.YellowString("[GetPost] Received: %v", in))
	id := in.GetPostID()

	// Get the post from the database
	post, err := s.sqlClient.GetPost(int(id))
	if err != nil {
		log.Fatal(color.RedString("[GetPost] DB error: %v", err))
	}

	response := &pb.GetPostResponse{Post: post}
	log.Print(color.GreenString("[GetPost] Reponse: %v", response))
	return response, nil
}

func main() {
	// Parse the flags
	flag.Parse()

	// Open the database
	s := &gRPCserver{}
	var err error
	s.sqlClient, err = NewSQLClient()
	if err != nil {
		log.Fatal(color.RedString("[Server] Error opening database: %v", err))
	}

	// Launch the server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *addr, *port))
	if err != nil {
		log.Fatal(color.RedString("[Server] Failed to listen: %v", err))
	}

	gs := grpc.NewServer()
	pb.RegisterRedditServer(gs, s)
	log.Print(color.GreenString("[Server] Listening at %v", lis.Addr()))
	if err := gs.Serve(lis); err != nil {
		log.Fatal(color.RedString("[Server] Failed to serve: %v", err))
	}
}
