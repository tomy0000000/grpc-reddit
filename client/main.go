package main

import (
	"flag"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	s := NewRedditAPIClient(*addr)

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
