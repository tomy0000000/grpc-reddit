package main

import (
	"flag"
	"log"

	"github.com/fatih/color"
)

var (
	addr = flag.String("addr", "localhost", "the address to connect to")
	port = flag.Int("port", 50051, "The server port")
)

func demoFunc(s RedditAPI) (string, error) {
	// Retrieve the post
	post, err := s.GetPost(1)
	if err != nil {
		log.Fatal(color.RedString("[GetPost] Error: %v", err))
		return "", err
	}
	log.Print(color.BlueString("[GetPost] Received: %v", post))

	// Retrieve most upvoted comments under the post
	comments, err := s.GetTopComments(post.Id, 10)
	if err != nil {
		log.Fatal(color.RedString("[GetTopComments] Error: %v", err))
		return "", err
	}
	log.Print(color.BlueString("[GetTopComments] Received: %v", comments))

	// Expand the most upvoted comment
	commentsOfComment, err := s.ExpandCommentBranch(comments[0].Id, 10)
	if err != nil {
		log.Fatal(color.RedString("[ExpandCommentBranch] Error: %v", err))
		return "", err
	}
	log.Print(color.BlueString("[ExpandCommentBranch] Received: %v", commentsOfComment))

	// Return the most upvoted reply under the most upvoted comment
	if len(commentsOfComment) == 0 {
		return "", nil
	}
	return commentsOfComment[0].GetContent(), nil
}

func main() {
	flag.Parse()
	s := NewRedditAPIClient(*addr, *port)

	log.Print(color.BlueString("[Demo] Start!"))
	result, err := demoFunc(s)
	if err != nil {
		log.Fatal(color.RedString("[Demo] Error: %v", err))
	}
	log.Print(color.BlueString("[Demo] Result: %v", result))
	log.Print(color.BlueString("[Demo] Done!"))

	s.runCreatePost()
	s.runVotePost()
	s.runGetPost()
	s.runCreateComment()
	s.runVoteComment()
	s.runGetComment()
	s.runGetTopComments()
	s.runExpandCommentBranch()
	s.runMonitorUpdates()

	s.close()
}
