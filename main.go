package main

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/sdehm/sdehm-blog-dynamic/data"
	"github.com/sdehm/sdehm-blog-dynamic/server"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                "https://c5607567d1f2bb6a06f0ac910baaa5fb@o4506211186573312.ingest.sentry.io/4506211186704384",
		TracesSampleRate:   1.0,
		EnableTracing:      true,
		ProfilesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	logger := log.New(log.Writer(), "server: ", log.Flags())
	// repo := data.NewDataMock()
	connectionString := os.Getenv("COCKROACH_CONNECTION")
	dataContext, cancel := context.WithCancel(ctx)
	defer cancel()

	repo, err := data.NewCockroachConnection(connectionString, dataContext)
	if err != nil {
		logger.Fatal("unable to create cockroach repo", err)
	}
	defer repo.Close()

	err = server.Start(":8080", logger, repo)
	if err != nil {
		logger.Fatal("unable to start server", err)
	}
}
