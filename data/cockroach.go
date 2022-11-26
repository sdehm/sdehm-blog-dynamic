package data

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/sdehm/sdehm-blog-dynamic/models"
)

type Cockroach struct {
	conn *pgx.Conn
}

func NewCockroachConnection() (*Cockroach, error) {
	// get connection string from environment variable
	// TODO: return the errors rather than log fatal
	config, err := pgx.ParseConfig(os.Getenv("COCKROACH_CONNECTION"))
	if err != nil {
		log.Fatal(" failed to parse config", err)
	}
	config.RuntimeParams["database"] = "blog"
	config.RuntimeParams["user"] = "blog"
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	// set database to "blog"
	// _, err = conn.Exec(context.Background(), "SET DATABASE = blog")
	// if err != nil {
	// 	log.Fatal("failed to set database", err)
	// }
	return &Cockroach{conn: conn}, nil
}

func (c *Cockroach) Close() error {
	return c.conn.Close(context.Background())
}

func (c *Cockroach) GetPost(path string) (*models.Post, error) {
	comments := []models.Comment{}
	rows, err := c.conn.Query(context.Background(), "SELECT id, author, body, created_at FROM comments WHERE post_path = $1", path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.Id, &comment.Author, &comment.Body, &comment.Timestamp)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return &models.Post{Comments: comments}, nil
}

func (c *Cockroach) AddComment(path string, author string, body string) (*models.Comment, error) {
	var comment models.Comment
	// id := uuid.New()
	timestamp := time.Now().UTC()
	c.addPostIfNotExists(path)
	err := c.conn.QueryRow(context.Background(), "INSERT INTO comments (post_path, author, body, created_at) VALUES ($1, $2, $3, $4) RETURNING id, created_at", path, author, body, timestamp).Scan(&comment.Id, &comment.Timestamp)
	if err != nil {
		return nil, err
	}
	comment.Author = author
	comment.Body = body
	return &comment, nil
}

func (c *Cockroach) addPostIfNotExists(path string) error {
	_, err := c.conn.Exec(context.Background(), "INSERT INTO posts (path) VALUES ($1) ON CONFLICT DO NOTHING", path)
	return err
}
