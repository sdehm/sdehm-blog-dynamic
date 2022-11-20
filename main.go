package main

import (
	"log"

	"github.com/sdehm/sdehm-blog-dynamic/data"
	"github.com/sdehm/sdehm-blog-dynamic/models"
	"github.com/sdehm/sdehm-blog-dynamic/server"
)

func main() {
	println("Hello, World!")
	logger := log.New(log.Writer(), "server: ", log.Flags())
	repo := data.NewDataMock()
	repo.AddComment("/", models.Comment{Id: 1, Author: "test", Body: "test"})
	server.Start(":8080", logger, repo)
}
