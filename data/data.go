package data

import (
	"github.com/sdehm/sdehm-blog-dynamic/models"
)

// Repo interface
type Repo interface {
	GetComment(id int) (*models.Comment, error)
	GetPost(path string) (models.Post, error)
	AddComment(string, string, string) (models.Comment, error)
	DeleteComment(id int) error
}
