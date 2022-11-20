package main

import "log"

func main() {
		println("Hello, World!")
		logger := log.New(log.Writer(), "server: ", log.Flags())
		repo := &dataMock{}
		start(":8080", logger, repo)
}
