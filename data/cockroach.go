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

func NewCockroachConnection(connectionString string, ctx context.Context) (*Cockroach, error) {
	// get connection string from environment variable
	config, err := pgx.ParseConfig(connectionString)
	// TODO: return the errors rather than log fatal
	if err != nil {
		log.Fatal(" failed to parse config", err)
	}
	config.RuntimeParams["database"] = "blog"
	config.RuntimeParams["user"] = "blog"
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	return &Cockroach{conn: conn, ctx: ctx}, nil
}

func (c *Cockroach) ReconnectIfClosed() error {
	if c.conn.IsClosed() {
		log.Default().Println("Connection is closed, reconnecting")
		conn, err := pgx.ConnectConfig(c.ctx, c.conn.Config())
		if err != nil {
			log.Fatal("failed to connect database", err)
		}
		c.conn = conn
	}
	return nil
}

func (c *Cockroach) Close() error {
	return c.conn.Close(c.ctx)
}

// TODO: may need to switch back to using two queries here if we ever want to
// get info frm the posts table since a post may have no comments
func (c *Cockroach) GetPost(path string) (*models.Post, error) {
	err := c.ReconnectIfClosed()
	if err != nil {
		return nil, err
	}
	sql := `SELECT p.id, p.path, c.id, c.author, c.body, c.created_at
			FROM posts p
			RIGHT JOIN comments c ON p.id = c.post_id
			WHERE p.path = $1
			ORDER BY c.created_at ASC`
	rows, err := c.conn.Query(c.ctx, sql, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	defer rows.Close()
	var post postDTO
	comments := []models.Comment{}
	for rows.Next() {
		var comment commentDTO
		err := rows.Scan(&post.Id, &post.Path, &comment.Id, &comment.Author, &comment.Body, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		comments = append(comments, models.Comment{
			Author:    comment.Author,
			Body:      comment.Body,
			Timestamp: comment.CreatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &models.Post{
		Path:     post.Path,
		Comments: comments,
	}, nil
}

func (c *Cockroach) AddComment(path string, author string, body string) (*models.Comment, error) {
	err := c.ReconnectIfClosed()
	if err != nil {
		return nil, err
	}
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
	err = c.conn.QueryRow(
		c.ctx,
		"INSERT INTO comments (post_id, author, body, created_at) VALUES ($1, $2, $3, $4) RETURNING author, body, created_at",
		post.Id, comment.Author, comment.Body, comment.CreatedAt).Scan(&comment.Author, &comment.Body, &comment.CreatedAt)
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
	err := c.ReconnectIfClosed()
	if err != nil {
		return nil, err
	}
	post := postDTO{}
	err = c.conn.QueryRow(c.ctx, "SELECT id, path FROM posts WHERE path = $1", path).Scan(&post.Id, &post.Path)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to add post: %w", err)
	}
	if err == nil {
		return &post, nil
	}
	err = c.conn.QueryRow(
		c.ctx,
		"INSERT INTO posts (path) VALUES ($1) RETURNING id, path",
		path).Scan(&post.Id, &post.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to add post: %w", err)
	}
	return &post, nil
}
