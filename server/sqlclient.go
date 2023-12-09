package main

import (
	"database/sql"
	"time"

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
		c.db.Exec("INSERT INTO post (title, content, subRedditID, videoURL, imageURL, authorID, score, state, publicationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			post.GetTitle(), post.GetContent(), post.GetSubReddit().GetId(),
			post.GetVideoURL(), post.GetImageURL(), post.GetAuthor().GetId(),
			post.GetScore(), post.GetState().Number(), time.Now(),
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
		&post.State, &post.PublicationDate,
	); err == sql.ErrNoRows {
		return nil, err
	}
	return post, nil
}

func (c *SQLClient) CreateComment(comment *pb.Comment) (int, error) {
	// Insert the comment into the database
	res, err :=
		c.db.Exec("INSERT INTO comment (content, authorID, score, state, publicationDate) VALUES (?, ?, ?, ?, ?)",
			comment.GetContent(), comment.GetAuthor().GetId(),
			comment.GetScore(), comment.GetState().Number(), time.Now(),
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

func (c *SQLClient) VoteComment(id int, upvote bool) (int, error) {
	// Increment/Decrement the score of the comment
	_, err := c.db.Exec("UPDATE comment SET score = score + (?) WHERE id = (?)", upvote, id)
	if err != nil {
		return -1, err
	}
	// Get the new score
	row := c.db.QueryRow("SELECT score FROM comment WHERE id = (?)", id)
	var newScore int
	if err := row.Scan(&newScore); err != nil {
		return -1, err
	}
	return newScore, nil
}

func (c *SQLClient) GetComment(id int) (*pb.Comment, error) {
	// Get the comment from the database
	row := c.db.QueryRow("SELECT * from comment WHERE id = (?)", id)
	comment := &pb.Comment{
		Author: &pb.User{},
	}
	if err := row.Scan(
		&comment.Id, &comment.Content, &comment.Author.Id, &comment.Score,
		&comment.State, &comment.PublicationDate, &comment.Parent, &comment.ParentID,
	); err == sql.ErrNoRows {
		return nil, err
	}
	return comment, nil
}

func (c *SQLClient) GetTopComments(postID int, quantity int) ([]*pb.Comment, error) {
	// Get the comment from the database
	rows, err := c.db.Query(
		"SELECT * from comment WHERE (parent = (?) AND parentID = (?)) ORDER BY score DESC LIMIT (?)",
		pb.CommentParent_POST, postID, quantity)
	if err != nil {
		return nil, err
	}
	comments := []*pb.Comment{}

	for rows.Next() {
		comment := &pb.Comment{
			Author: &pb.User{},
		}
		if err := rows.Scan(
			&comment.Id, &comment.Content, &comment.Author.Id, &comment.Score,
			&comment.State, &comment.PublicationDate, &comment.Parent, &comment.ParentID,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (c *SQLClient) ExpandCommentBranch(id int) ([]*pb.Comment, error) {
	// Get the comment from the database
	rows, err := c.db.Query(
		"SELECT * from comment WHERE (parent = (?) AND parentID = (?)) ORDER BY score DESC",
		pb.CommentParent_COMMENT, id)
	if err != nil {
		return nil, err
	}
	comments := []*pb.Comment{}

	for rows.Next() {
		comment := &pb.Comment{
			Author: &pb.User{},
		}
		if err := rows.Scan(
			&comment.Id, &comment.Content, &comment.Author.Id, &comment.Score,
			&comment.State, &comment.PublicationDate, &comment.Parent, &comment.ParentID,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
