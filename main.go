package main

import (
	"log"

	"github.com/sdehm/sdehm-blog-dynamic/data"
	"github.com/sdehm/sdehm-blog-dynamic/server"
)

func main() {
	logger := log.New(log.Writer(), "server: ", log.Flags())
	// repo := data.NewDataMock()
	repo, err := data.NewCockroachConnection()
	if err != nil {
		logger.Fatal("unable to create cockroach repo", err)
	}
	defer repo.Close()
	server.Start(":8080", logger, repo)
}
