package main

import (
	"log"
	"os"

	"github.com/oss-health/background-worker/fetcher"
	"github.com/oss-health/background-worker/parser"
	"github.com/oss-health/background-worker/pipeline"
	"github.com/oss-health/background-worker/resolver"
	"github.com/oss-health/background-worker/storage"
	"golang.org/x/time/rate"

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

	npmResolver := resolver.NewRateLimitedResolver(
		&resolver.NPMResolver{},
		rate.Limit(10),
		5,
	)
	npmResolver.StartLogging("npm")
	pypiResolver := resolver.NewRateLimitedResolver(
		&resolver.PyPIResolver{},
		rate.Limit(10),
		5,
	)
	pypiResolver.StartLogging("pypi")

	resolverChain := resolver.NewResolverChain(
		npmResolver,
		pypiResolver,
	)

	p := pipeline.NewPipeline(
		fetcher.NewFetcherChain(),
		parser.NewParserRegistry(),
		store,
		resolverChain,
	)

	if err := p.Run(); err != nil {
		log.Fatalf("failed to run pipeline: %v", err)

	}
}
