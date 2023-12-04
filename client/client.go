package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/tomy0000000/grpc-reddit/reddit/reddit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "Tomy"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRedditClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreatePost(ctx, &pb.CreatePostRequest{Post: &pb.Post{Title: "Hello", Content: "World", SubReddit: &pb.SubReddit{Id: 1}, Author: &pb.User{Id: 1}}})
	if err != nil {
		log.Fatalf("could not create post: %v", err)
	}
	log.Printf("Post Created: %v", r.GetPost())

	r2, err := c.GetPost(ctx, &pb.GetPostRequest{PostID: int32(1)})
	if err != nil {
		log.Fatalf("could not get post: %v", err)
	}
	log.Printf("Post Retrieved: %v", r2.GetPost())
}
