package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/fatih/color"
	pb "github.com/tomy0000000/grpc-reddit/reddit/reddit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	TIMEOUT = 1000 * time.Second
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func runCreatePost(c pb.RedditClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.CreatePostRequest{
		Post: &pb.Post{
			Title:     "Hello",
			Content:   "World",
			SubReddit: &pb.SubReddit{Id: 1},
			Author:    &pb.User{Id: 1},
		},
	}
	log.Print(color.YellowString("[CreatePost] Sending: %v", request))

	response, err := c.CreatePost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[CreatePost] Error: %v", err))
	}

	log.Print(color.GreenString("[CreatePost] Received: %v", response))
}

func runVotePost(c pb.RedditClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.VotePostRequest{PostID: int32(1), Upvote: true}
	log.Print(color.YellowString("[VotePost] Sending: %v", request))

	response, err := c.VotePost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[VotePost] Error: %v", err))
	}

	log.Print(color.GreenString("[VotePost] Received: %v", response))
}

func runGetPost(c pb.RedditClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.GetPostRequest{PostID: int32(1)}
	log.Print(color.YellowString("[GetPost] Sending: %v", request))

	response, err := c.GetPost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[GetPost] Error: %v", err))
	}
	log.Print(color.GreenString("[GetPost] Received: %v", response))
}

func runCreateComment(c pb.RedditClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.CreateCommentRequest{
		Comment: &pb.Comment{
			Content: "Hello",
			Author:  &pb.User{Id: 1},
		},
	}
	log.Print(color.YellowString("[CreateComment] Sending: %v", request))

	response, err := c.CreateComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[CreateComment] Error: %v", err))
	}
	log.Print(color.GreenString("[CreateComment] Received: %v", response))
}

func runVoteComment(c pb.RedditClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.VoteCommentRequest{CommentID: int32(1), Upvote: true}
	log.Print(color.YellowString("[VoteComment] Sending: %v", request))

	response, err := c.VoteComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[VoteComment] Error: %v", err))
	}
	log.Print(color.GreenString("[VoteComment] Received: %v", response))
}

func runGetComment(c pb.RedditClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.GetCommentRequest{CommentID: int32(1)}
	log.Print(color.YellowString("[GetComment] Sending: %v", request))

	response, err := c.GetComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[GetComment] Error: %v", err))
	}
	log.Print(color.GreenString("[GetComment] Received: %v", response))
}

func runGetTopComments(c pb.RedditClient) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	requests := &pb.GetTopCommentsRequest{PostID: int32(2), Quantity: int32(10)}
	log.Print(color.YellowString("[GetTopComments] Sending: %v", requests))

	response, err := c.GetTopComments(ctx, requests)
	if err != nil {
		log.Fatal(color.RedString("[GetTopComments] Error: %v", err))
	}
	log.Print(color.GreenString("[GetTopComments] Received: %v", response))
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(color.RedString("did not connect: %v", err))
	}
	defer conn.Close()
	c := pb.NewRedditClient(conn)

	runCreatePost(c)
	runVotePost(c)
	runGetPost(c)
	runCreateComment(c)
	runVoteComment(c)
	runGetComment(c)
	runGetTopComments(c)
}
