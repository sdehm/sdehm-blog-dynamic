package main

import (
	"log"
	"os"

	"github.com/sdehm/sdehm-blog-dynamic/data"
	"github.com/sdehm/sdehm-blog-dynamic/server"
	"golang.org/x/net/context"
)

func main() {
	logger := log.New(log.Writer(), "server: ", log.Flags())
	// repo := data.NewDataMock()
	connectionString := os.Getenv("COCKROACH_CONNECTION")
	dataContext := context.Background()
	dataContext, cancel := context.WithCancel(dataContext)
	defer cancel()

	repo, err := data.NewCockroachConnection(connectionString, dataContext)
	if err != nil {
		logger.Fatal("unable to create cockroach repo", err)
	}
	defer repo.Close()

	server.Start(":8080", logger, repo)
}
