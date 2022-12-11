package data

import (
	"github.com/sdehm/sdehm-blog-dynamic/models"
)

type Repo interface {
	GetPost(path string) (*models.Post, error)
	AddComment(string, string, string) (*models.Comment, error)
}
