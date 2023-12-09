package main

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	pb "github.com/tomy0000000/grpc-reddit/reddit/reddit"
)

const (
	db_file string = "data/reddit.db"
)

type SQLClient struct {
	db *sql.DB
}

func NewSQLClient() (*SQLClient, error) {
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		return nil, err
	}
	return &SQLClient{db: db}, nil
}

func (c *SQLClient) CreatePost(post *pb.Post) (int, error) {
	// Insert the post into the database
	res, err :=
		c.db.Exec("INSERT INTO post (title, content, subRedditID, videoURL, imageURL, authorID, score, state, publicationDate, comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			post.GetTitle(), post.GetContent(), post.GetSubReddit().GetId(),
			post.GetVideoURL(), post.GetImageURL(), post.GetAuthor().GetId(),
			post.GetScore(), post.GetState().Number(), post.GetPublicationDate(),
			strings.Join(post.GetComments(), ","),
		)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func (c *SQLClient) VotePost(id int, upvote bool) (int, error) {
	// Increment/Decrement the score of the post
	_, err := c.db.Exec("UPDATE post SET score = score + (?) WHERE id = (?)", upvote, id)
	if err != nil {
		return -1, err
	}
	// Get the new score
	row := c.db.QueryRow("SELECT score FROM post WHERE id = (?)", id)
	var newScore int
	if err := row.Scan(&newScore); err != nil {
		return -1, err
	}
	return newScore, nil
}

func (c *SQLClient) GetPost(id int) (*pb.Post, error) {
	// Get the post from the database
	row := c.db.QueryRow("SELECT * from post WHERE id = (?)", id)
	post := &pb.Post{
		SubReddit: &pb.SubReddit{},
		Author:    &pb.User{},
	}
	if err := row.Scan(
		&post.Id, &post.Title, &post.Content, &post.SubReddit.Id,
		&post.VideoURL, &post.ImageURL, &post.Author.Id, &post.Score,
		&post.State, &post.PublicationDate, &post.Comments,
	); err == sql.ErrNoRows {
		return nil, err
	}
	return post, nil
}
