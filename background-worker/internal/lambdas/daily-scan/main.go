package main

import (
	"context"
	"log"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repository"
	"github.com/oss-health/background-worker/internal/dependency"
)

func main() {
	connectionString := "postgres://dev-user:password@localhost:5432/dev_db"
	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	repo := repository.NewRepositoryRepository(db)
	dependencyRepository := dependency.NewPostgresRepository(db)
	cloner := &repository.GitCloner{}
	extractor := &repository.DependencyExtractor{}
	service := repository.NewRepositoryService(repo, dependencyRepository, cloner, extractor)
	// service.RunDailyScan(ctx, 4, 3)
	service.RunDailyScan(ctx, 6, 9)
}
