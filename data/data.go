package data

import (
	"time"

	"github.com/sdehm/sdehm-blog-dynamic/models"
)

// Repo interface
type Repo interface {
	GetComment(id int) (*models.Comment, error)
	GetPost(path string) (models.Post, error)
	AddComment(string, string, string) (models.Comment, error)
	DeleteComment(id int) error
}

// mock implementation of the data interface for testing
type DataMock struct {
	posts map[string]models.Post
}

func NewDataMock() *DataMock {
	return &DataMock{
		posts: map[string]models.Post{},
	}
}

func (d *DataMock) GetComment(id int) (*models.Comment, error) {
	for _, p := range d.posts {
		for _, c := range p.Comments {
			if c.Id == id {
				return &c, nil
			}
		}
	}
	return nil, nil
}

func (d *DataMock) GetPost(path string) (models.Post, error) {
	return d.posts[path], nil
}

func (d *DataMock) AddComment(p string, author string, body string) (models.Comment, error) {
	c := models.Comment{
		Id:       len(d.posts[p].Comments) + 1,
		Author:  author,
		Body: body,
		Timestamp: 	time.Now().UTC(),
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

func (d *DataMock) DeleteComment(id int) error {
	for _, p := range d.posts {
		for i, c := range p.Comments {
			if c.Id == id {
				p.Comments = append(p.Comments[:i], p.Comments[i+1:]...)
				return nil
			}
		}
	}
	return nil
}
