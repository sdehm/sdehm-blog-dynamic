package main

import (
	"log"
	"time"

	"github.com/sdehm/sdehm-blog-dynamic/data"
	"github.com/sdehm/sdehm-blog-dynamic/models"
	"github.com/sdehm/sdehm-blog-dynamic/server"
)

func main() {
	println("Hello, World!")
	logger := log.New(log.Writer(), "server: ", log.Flags())
	repo := data.NewDataMock()
	repo.AddComment("/", models.Comment{Id: 1, Author: "Author", Body: "Comment body goes here.", Timestamp: time.Now()})
	repo.AddComment("/", models.Comment{Id: 2, Author: "Other Author", Body: "Comment body for the other comment is here.", Timestamp: time.Now()})
	server.Start(":8080", logger, repo)
}
