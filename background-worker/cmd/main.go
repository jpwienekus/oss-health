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

	repoUrl := "https://github.com/jpwienekus/oss-health.git"
	if err := p.Run(repoUrl); err != nil {
		log.Fatalf("failed to run pipeline: %v", err)

	}

	// repoUrl := "https://github.com/jpwienekus/oss-health.git"
	// dest, err := os.MkdirTemp("", "repo-clone-*")
	// if err != nil {
	// 	log.Fatalf("failed to create temp directory: %v", err)
	// }
	//
	// fetcherChain := fetcher.NewFetcherChain()
	// if err := fetcherChain.Fetch(repoUrl, dest); err != nil {
	// 	log.Fatalf("failed to fetch repository: %v", err)
	// }
	//
	// if removeErr := os.RemoveAll(dest); removeErr != nil {
	// 	log.Fatalf("failed to remove temporary directory: %v", removeErr)
	// }
	//
	// fmt.Printf("Repository cloned to: %s\n", dest)
}
