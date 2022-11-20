package models

import (
	"time"
)

type Comment struct {
	Id        int
	Author    string
	Body      string
	Timestamp time.Time
}

type Post struct {
	Path     string
	Comments []Comment
}
