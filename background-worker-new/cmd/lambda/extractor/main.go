package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies"
	"github.com/oss-health/background-worker/internal/repo"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	dependencyService := dependencies.NewDependencyService(
		dependencies.NewDependencyRepository(db),
	)

	service := repo.NewRepositoryService(
		repo.NewRepoRepository(db),
		dependencyService,
		parsers.NewParserRegistry(),

		// Add other dependencies here if needed
	)

	service.RunDailyScan(ctx, 0, 0)
}
