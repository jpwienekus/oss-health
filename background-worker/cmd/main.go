package main

import (
	"log"
	"os"

	"github.com/oss-health/background-worker/fetcher"
	"github.com/oss-health/background-worker/parser"
	"github.com/oss-health/background-worker/pipeline"
	"github.com/oss-health/background-worker/storage"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("DATABASE_URL")

	store, err := storage.NewPostgresStore(connectionString)
	if err != nil {
		log.Fatalf("failed to create postgres store: %v", err)
	}

	p := pipeline.NewPipeline(
		fetcher.NewFetcherChain(),
		parser.NewParserRegistry(),
		store,
	)

	if err := p.Run(); err != nil {
		log.Fatalf("failed to run pipeline: %v", err)

	}
}
