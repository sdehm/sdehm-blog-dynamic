package models

import (
	"time"
)

type Comment struct {
	Author    string
	Body      string
	Timestamp time.Time
}

type Post struct {
	Path     string
	Comments []Comment
}
