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

type RedditUser = pb.User
type RedditSubReddit = pb.SubReddit
type RedditPost = pb.Post
type RedditComment = pb.Comment

type RedditAPI interface {
	GetPost(PostID int32) (*RedditPost, error)
	GetTopComments(PostID int32, Quantity int32) ([]*RedditComment, error)
	ExpandCommentBranch(commentId int32, Quantity int32) ([]*RedditComment, error)
}

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

func (s *RedditAPIClient) CreatePost(title string, content string, subRedditID int32, authorID int32) (*RedditPost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.CreatePostRequest{
		Post: &RedditPost{
			Title:     title,
			Content:   content,
			SubReddit: &RedditSubReddit{Id: subRedditID},
			Author:    &RedditUser{Id: authorID},
		},
	}
	log.Print(color.YellowString("[CreatePost] Sending: %v", request))

	response, err := s._client.CreatePost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[CreatePost] Error: %v", err))
		return nil, err
	}

	log.Print(color.GreenString("[CreatePost] Received: %v", response))
	return response.Post, nil
}

func (s *RedditAPIClient) VotePost(postID int32, upvote bool) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.VotePostRequest{PostID: postID, Upvote: upvote}
	log.Print(color.YellowString("[VotePost] Sending: %v", request))

	response, err := s._client.VotePost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[VotePost] Error: %v", err))
		return -1, err
	}

	log.Print(color.GreenString("[VotePost] Received: %v", response))
	return response.Score, nil
}

func (s *RedditAPIClient) GetPost(postID int32) (*RedditPost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.GetPostRequest{PostID: postID}
	log.Print(color.YellowString("[GetPost] Sending: %v", request))

	response, err := s._client.GetPost(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[GetPost] Error: %v", err))
		return nil, err
	}
	log.Print(color.GreenString("[GetPost] Received: %v", response))
	return response.Post, nil
}

func (s *RedditAPIClient) CreateComment(authorID int32, content string) (*RedditComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.CreateCommentRequest{
		Comment: &RedditComment{
			Content:  content,
			Author:   &RedditUser{Id: authorID},
			Parent:   pb.ContentType_POST,
			ParentID: 1,
		},
	}
	log.Print(color.YellowString("[CreateComment] Sending: %v", request))

	response, err := s._client.CreateComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[CreateComment] Error: %v", err))
		return nil, err
	}
	log.Print(color.GreenString("[CreateComment] Received: %v", response))
	return response.Comment, nil
}

func (s *RedditAPIClient) VoteComment(commentID int32, upvote bool) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.VoteCommentRequest{CommentID: commentID, Upvote: upvote}
	log.Print(color.YellowString("[VoteComment] Sending: %v", request))

	response, err := s._client.VoteComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[VoteComment] Error: %v", err))
		return -1, err
	}
	log.Print(color.GreenString("[VoteComment] Received: %v", response))
	return response.Score, nil
}

func (s *RedditAPIClient) GetComment(commentID int32) (*RedditComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	request := &pb.GetCommentRequest{CommentID: commentID}
	log.Print(color.YellowString("[GetComment] Sending: %v", request))

	response, err := s._client.GetComment(ctx, request)
	if err != nil {
		log.Fatal(color.RedString("[GetComment] Error: %v", err))
		return nil, err
	}
	log.Print(color.GreenString("[GetComment] Received: %v", response))
	return response.Comment, nil
}

func (s *RedditAPIClient) GetTopComments(postID int32, quantity int32) ([]*RedditComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	requests := &pb.GetTopCommentsRequest{PostID: postID, Quantity: quantity}
	log.Print(color.YellowString("[GetTopComments] Sending: %v", requests))

	response, err := s._client.GetTopComments(ctx, requests)
	if err != nil {
		log.Fatal(color.RedString("[GetTopComments] Error: %v", err))
		return nil, err
	}
	log.Print(color.GreenString("[GetTopComments] Received: %v", response))
	return response.Comments, nil
}

func (s *RedditAPIClient) ExpandCommentBranch(commentID int32, quantity int32) ([]*RedditComment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	requests := &pb.ExpandCommentBranchRequest{CommentID: commentID, Quantity: quantity}
	log.Print(color.YellowString("[ExpandCommentBranch] Sending: %v", requests))

	response, err := s._client.ExpandCommentBranch(ctx, requests)
	if err != nil {
		log.Fatal(color.RedString("[ExpandCommentBranch] Error: %v", err))
		return nil, err
	}
	log.Print(color.GreenString("[ExpandCommentBranch] Received: %v", response))
	return response.Comments, nil
}

func (s *RedditAPIClient) runCreatePost() {
	s.CreatePost("Hello", "World", 1, 1)
}

func (s *RedditAPIClient) runVotePost() {
	s.VotePost(1, true)
}

func (s *RedditAPIClient) runGetPost() {
	s.GetPost(1)
}

func (s *RedditAPIClient) runCreateComment() {
	s.CreateComment(1, "Hello World")
}

func (s *RedditAPIClient) runVoteComment() {
	s.VoteComment(1, true)
}

func (s *RedditAPIClient) runGetComment() {
	s.GetComment(1)
}

func (s *RedditAPIClient) runGetTopComments() {
	s.GetTopComments(2, 10)
}

func (s *RedditAPIClient) runExpandCommentBranch() {
	s.ExpandCommentBranch(1, 10)
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
