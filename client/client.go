package main

import (
	"context"
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

type RedditAPIClient struct {
	_client pb.RedditClient
	_conn   *grpc.ClientConn
}

func NewRedditAPIClient(addr string) *RedditAPIClient {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(color.RedString("did not connect: %v", err))
	}

	return &RedditAPIClient{
		_client: pb.NewRedditClient(conn),
		_conn:   conn,
	}
}

func (s *RedditAPIClient) close() {
	s._conn.Close()
}

func (s *RedditAPIClient) runCreatePost() {
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

	response, err := s._client.CreatePost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[CreatePost] Error: %v", err))
	}

	log.Print(color.GreenString("[CreatePost] Received: %v", response))
}

func (s *RedditAPIClient) runVotePost() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.VotePostRequest{PostID: int32(1), Upvote: true}
	log.Print(color.YellowString("[VotePost] Sending: %v", request))

	response, err := s._client.VotePost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[VotePost] Error: %v", err))
	}

	log.Print(color.GreenString("[VotePost] Received: %v", response))
}

func (s *RedditAPIClient) runGetPost() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.GetPostRequest{PostID: int32(1)}
	log.Print(color.YellowString("[GetPost] Sending: %v", request))

	response, err := s._client.GetPost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[GetPost] Error: %v", err))
	}
	log.Print(color.GreenString("[GetPost] Received: %v", response))
}

func (s *RedditAPIClient) runCreateComment() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.CreateCommentRequest{
		Comment: &pb.Comment{
			Content: "Hello",
			Author:  &pb.User{Id: 1},
		},
	}
	log.Print(color.YellowString("[CreateComment] Sending: %v", request))

	response, err := s._client.CreateComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[CreateComment] Error: %v", err))
	}
	log.Print(color.GreenString("[CreateComment] Received: %v", response))
}

func (s *RedditAPIClient) runVoteComment() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.VoteCommentRequest{CommentID: int32(1), Upvote: true}
	log.Print(color.YellowString("[VoteComment] Sending: %v", request))

	response, err := s._client.VoteComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[VoteComment] Error: %v", err))
	}
	log.Print(color.GreenString("[VoteComment] Received: %v", response))
}

func (s *RedditAPIClient) runGetComment() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.GetCommentRequest{CommentID: int32(1)}
	log.Print(color.YellowString("[GetComment] Sending: %v", request))

	response, err := s._client.GetComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[GetComment] Error: %v", err))
	}
	log.Print(color.GreenString("[GetComment] Received: %v", response))
}

func (s *RedditAPIClient) runGetTopComments() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	requests := &pb.GetTopCommentsRequest{PostID: int32(2), Quantity: int32(10)}
	log.Print(color.YellowString("[GetTopComments] Sending: %v", requests))

	response, err := s._client.GetTopComments(ctx, requests)
	if err != nil {
		log.Fatal(color.RedString("[GetTopComments] Error: %v", err))
	}
	log.Print(color.GreenString("[GetTopComments] Received: %v", response))
}

func (s *RedditAPIClient) runExpandCommentBranch() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	requests := &pb.ExpandCommentBranchRequest{CommentID: int32(1)}
	log.Print(color.YellowString("[ExpandCommentBranch] Sending: %v", requests))

	response, err := s._client.ExpandCommentBranch(ctx, requests)
	if err != nil {
		log.Fatal(color.RedString("[ExpandCommentBranch] Error: %v", err))
	}
	log.Print(color.GreenString("[ExpandCommentBranch] Received: %v", response))
}

func (s *RedditAPIClient) runMonitorUpdates() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	// Create a stream
	stream, err := s._client.MonitorUpdates(ctx)
	if err != nil {
		log.Fatal(color.RedString("[MonitorUpdates] Error: %v", err))
	}

	// Routine to receive responses
	waitc := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Print(color.GreenString("[MonitorUpdates] Received: %v", response))
		}
	}()

	// Send a initial monitor request
	requests := &pb.MonitorUpdatesRequest{
		ContentType: pb.ContentType_POST, ContentID: int32(1),
	}
	log.Print(color.YellowString("[MonitorUpdates] Sending: %v", requests))
	if err := stream.Send(requests); err != nil {
		log.Fatal(color.RedString("[MonitorUpdates] Error: %v", err))
	}

	// Wait for 10 seconds
	time.Sleep(10 * time.Second)

	// Send a second monitor request
	requests = &pb.MonitorUpdatesRequest{
		ContentType: pb.ContentType_COMMENT, ContentID: int32(1),
	}
	log.Print(color.YellowString("[MonitorUpdates] Sending: %v", requests))
	if err := stream.Send(requests); err != nil {
		log.Fatal(color.RedString("[MonitorUpdates] Error: %v", err))
	}

	// Close the stream after 10 seconds
	time.Sleep(10 * time.Second)
	stream.CloseSend()
}
