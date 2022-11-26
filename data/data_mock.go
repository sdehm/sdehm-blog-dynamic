package data

import (
	"time"

	"github.com/sdehm/sdehm-blog-dynamic/models"
)

// mock implementation of the data interface for testing
type DataMock struct {
	posts map[string]models.Post
}

func NewDataMock() *DataMock {
	return &DataMock{
		posts: map[string]models.Post{},
	}
}

func (d *DataMock) GetPost(path string) (models.Post, error) {
	return d.posts[path], nil
}

func (d *DataMock) AddComment(p string, author string, body string) (models.Comment, error) {
	c := models.Comment{
		Author:    author,
		Body:      body,
		Timestamp: time.Now().UTC(),
	}
	post, ok := d.posts[p]
	if !ok {
		d.posts[p] = models.Post{
			Path:     p,
			Comments: []models.Comment{c},
		}
	}
	post.Comments = append(post.Comments, c)
	d.posts[p] = post
	return c, nil
}
