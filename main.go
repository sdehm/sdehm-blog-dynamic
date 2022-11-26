package main

import (
	"log"
	"os"

	"github.com/sdehm/sdehm-blog-dynamic/data"
	"github.com/sdehm/sdehm-blog-dynamic/server"
)

func main() {
	logger := log.New(log.Writer(), "server: ", log.Flags())
	// repo := data.NewDataMock()
	connectionString := os.Getenv("COCKROACH_CONNECTION")
	repo, err := data.NewCockroachConnection(connectionString)
	if err != nil {
		logger.Fatal("unable to create cockroach repo", err)
	}
	defer repo.Close()
	server.Start(":8080", logger, repo)
}
