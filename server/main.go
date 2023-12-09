package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	pb "github.com/tomy0000000/grpc-reddit/reddit/reddit"
	"google.golang.org/grpc"
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

func (s *gRPCserver) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	log.Print(color.YellowString("[CreateComment] Received: %v", in))

	// Insert the comment into the database
	id, err := s.sqlClient.CreateComment(in.GetComment())
	if err != nil {
		log.Fatal(color.RedString("[CreateComment] DB error: %v", err))
	}

	// Get the comment from the database
	comment, err := s.sqlClient.GetComment(id)
	if err != nil {
		log.Fatal(color.RedString("[CreateComment] DB error: %v", err))
	}

	response := &pb.CreateCommentResponse{Comment: comment}
	log.Print(color.GreenString("[CreateComment] Reponse: %v", response))
	return response, nil
}

func (s *gRPCserver) VoteComment(ctx context.Context, in *pb.VoteCommentRequest) (*pb.VoteCommentResponse, error) {
	log.Print(color.YellowString("[VoteComment] Received: %v", in))

	// Increment/Decrement the score of the post
	newScore, err := s.sqlClient.VoteComment(int(in.GetCommentID()), in.GetUpvote())
	if err != nil {
		log.Fatal(color.RedString("[VoteComment] DB error: %v", err))
	}

	response := &pb.VoteCommentResponse{Score: int32(newScore)}
	log.Print(color.GreenString("[VoteComment] Reponse: %v", response))
	return response, nil
}

func (s *gRPCserver) GetComment(ctx context.Context, in *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	log.Print(color.YellowString("[GetComment] Received: %v", in))
	id := in.GetCommentID()

	// Get the comment from the database
	comment, err := s.sqlClient.GetComment(int(id))
	if err != nil {
		log.Fatal(color.RedString("[GetComment] DB error: %v", err))
	}

	response := &pb.GetCommentResponse{Comment: comment}
	log.Print(color.GreenString("[GetComment] Reponse: %v", response))
	return response, nil
}

func (s *gRPCserver) GetTopComments(ctx context.Context, in *pb.GetTopCommentsRequest) (*pb.GetTopCommentsResponse, error) {
	log.Print(color.YellowString("[GetTopComments] Received: %v", in))

	// Get the comments from the database
	comments, err := s.sqlClient.GetTopComments(int(in.GetPostID()), int(in.GetQuantity()))
	if err != nil {
		log.Fatal(color.RedString("[GetTopComments] DB error: %v", err))
	}

	response := &pb.GetTopCommentsResponse{Comments: comments}
	log.Print(color.GreenString("[GetTopComments] Reponse: %v", response))
	return response, nil
}

func (s *gRPCserver) ExpandCommentBranch(ctx context.Context, in *pb.ExpandCommentBranchRequest) (*pb.ExpandCommentBranchResponse, error) {
	log.Print(color.YellowString("[ExpandCommentBranch] Received: %v", in))

	// Get the comments from the database
	comments, err := s.sqlClient.ExpandCommentBranch(int(in.GetCommentID()))
	if err != nil {
		log.Fatal(color.RedString("[ExpandCommentBranch] DB error: %v", err))
	}

	response := &pb.ExpandCommentBranchResponse{Comments: comments}
	log.Print(color.GreenString("[ExpandCommentBranch] Reponse: %v", response))
	return response, nil
}

func (s *gRPCserver) MonitorUpdates(stream pb.Reddit_MonitorUpdatesServer) error {
	monitorPostList := []int{}
	monitorCommentList := []int{}

	// Process client requests to add content to the list of monitored contents
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Print(color.RedString("[MonitorUpdates] Error: %v", err))
				return
			}
			log.Print(color.YellowString("[MonitorUpdates] Received: %v", in))

			// Add the content to the list of monitored contents
			switch in.GetContentType() {
			case pb.ContentType_POST:
				monitorPostList = append(monitorPostList, int(in.GetContentID()))
			case pb.ContentType_COMMENT:
				monitorCommentList = append(monitorCommentList, int(in.GetContentID()))
			}
		}
	}()

	// Send the updates
	for {
		// Send the updates for the posts
		for _, postID := range monitorPostList {
			post, err := s.sqlClient.GetPost(postID)
			if err != nil {
				log.Fatal(color.RedString("[MonitorUpdates] DB error: %v", err))
			}

			// Send the updates
			response := &pb.MonitorUpdatesResponse{
				ContentType: pb.ContentType_POST,
				ContentID:   int32(postID),
				Score:       post.Score,
			}
			log.Print(color.GreenString("[MonitorUpdates] Response: %v", response))
			err = stream.Send(response)
			if err != nil {
				log.Print(color.RedString("[MonitorUpdates] Error: %v", err))
				return err
			}
		}

		// Send the updates for the comments
		for _, commentID := range monitorCommentList {
			comment, err := s.sqlClient.GetComment(commentID)
			if err != nil {
				log.Fatal(color.RedString("[MonitorUpdates] DB error: %v", err))
			}

			// Send the updates
			response := &pb.MonitorUpdatesResponse{
				ContentType: pb.ContentType_COMMENT,
				ContentID:   int32(commentID),
				Score:       comment.Score,
			}
			log.Print(color.GreenString("[MonitorUpdates] Response: %v", response))
			err = stream.Send(response)
			if err != nil {
				log.Print(color.RedString("[MonitorUpdates] Error: %v", err))
				return err
			}
		}

		// Wait for 2 seconds before sending the updates again
		time.Sleep(2 * time.Second)
	}
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
