package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repo"
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

	service := repo.NewRepositoryService(
		repo.NewRepoRepository(db),
		// Add other dependencies here if needed
	)

	service.RunDailyScan(ctx, 0, 0)
}
