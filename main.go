package main

import (
	"log"

	"github.com/sdehm/sdehm-blog-dynamic/data"
	"github.com/sdehm/sdehm-blog-dynamic/server"
)

func main() {
	logger := log.New(log.Writer(), "server: ", log.Flags())
	repo := data.NewDataMock()
	server.Start(":8080", logger, repo)
}
