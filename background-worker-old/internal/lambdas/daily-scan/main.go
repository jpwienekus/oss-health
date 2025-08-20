package main

import (
	"context"
	"log"
	"os"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/repository"
	"github.com/oss-health/background-worker/internal/repository/parsers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	repo := repository.NewRepositoryRepository(db)
	dependencyRepository := dependency.NewPostgresRepository(db)
	cloner := &repository.GitCloner{}
	provider := &parsers.ParserProviderImpl{}
	extractor := &repository.DependencyExtractor{
		Provider: provider,
	}
	service := repository.NewRepositoryService(repo, dependencyRepository, cloner, extractor)
	// service.RunDailyScan(ctx, 4, 3)
	if err := service.RunDailyScan(ctx, 0, 0); err != nil {
		log.Printf("error running daily scan: %v", err)
	}
}
