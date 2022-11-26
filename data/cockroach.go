package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sdehm/sdehm-blog-dynamic/models"
)

type Cockroach struct {
	conn *pgx.Conn
	ctx  context.Context
}

type commentDTO struct {
	Id        uuid.UUID
	Author    string
	Body      string
	CreatedAt time.Time
}

type postDTO struct {
	Id   uuid.UUID
	Path string
}

func NewCockroachConnection(connectionString string) (*Cockroach, error) {
	context := context.TODO()
	// get connection string from environment variable
	// TODO: return the errors rather than log fatal
	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		log.Fatal(" failed to parse config", err)
	}
	config.RuntimeParams["database"] = "blog"
	config.RuntimeParams["user"] = "blog"
	conn, err := pgx.ConnectConfig(context, config)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	// set database to "blog"
	// _, err = conn.Exec(c.ctx, "SET DATABASE = blog")
	// if err != nil {
	// 	log.Fatal("failed to set database", err)
	// }
	return &Cockroach{conn: conn, ctx: context}, nil
}

func (c *Cockroach) Close() error {
	return c.conn.Close(c.ctx)
}

func (c *Cockroach) GetPost(path string) (*models.Post, error) {
	post := postDTO{}
	err := c.conn.QueryRow(c.ctx, "SELECT id, path FROM posts WHERE path = $1", path).Scan(&post.Id, &post.Path)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	comments := []models.Comment{}
	// TODO: use a join to get the comments
	rows, err := c.conn.Query(c.ctx, "SELECT id, author, body, created_at FROM comments WHERE post_id = $1", post.Id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		comment := commentDTO{}
		err = rows.Scan(&comment.Id, &comment.Author, &comment.Body, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, models.Comment{
			Author:    comment.Author,
			Body:      comment.Body,
			Timestamp: comment.CreatedAt,
		})
	}
	return &models.Post{
		Path:     post.Path,
		Comments: comments,
	}, nil
}

func (c *Cockroach) AddComment(path string, author string, body string) (*models.Comment, error) {
	// TODO: use a transaction
	comment := commentDTO{
		Author:    author,
		Body:      body,
		CreatedAt: time.Now().UTC(),
	}
	post, err := c.getOrAddPost(path)
	if err != nil {
		return nil, err
	}
	err = c.conn.QueryRow(c.ctx, "INSERT INTO comments (post_id, author, body, created_at) VALUES ($1, $2, $3, $4) RETURNING author, body, created_at", post.Id, comment.Author, comment.Body, comment.CreatedAt).Scan(&comment.Author, &comment.Body, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}
	return &models.Comment{
		Author:    comment.Author,
		Body:      comment.Body,
		Timestamp: comment.CreatedAt,
	}, nil
}

func (c *Cockroach) getOrAddPost(path string) (*postDTO, error) {
	post := postDTO{}
	err := c.conn.QueryRow(c.ctx, "SELECT id, path FROM posts WHERE path = $1", path).Scan(&post.Id, &post.Path)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to add post: %w", err)
	}
	if err == nil {
		return &post, nil
	}
	err = c.conn.QueryRow(c.ctx, "INSERT INTO posts (path) VALUES ($1) RETURNING id, path", path).Scan(&post.Id, &post.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to add post: %w", err)
	}
	return &post, nil
}
