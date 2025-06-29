package main

import (
	"context"
	"log"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/repository"
	"github.com/oss-health/background-worker/internal/repository/parsers"
)

func main() {
	// connectionString := "postgres://dev-user:password@localhost:5432/dev_db"
	connectionString := "postgresql://postgres.gfpivacysduostopkekw:4-dzBCK8Ptyg.FTukiBB@aws-0-eu-central-1.pooler.supabase.com:5432/postgres"
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
	service.RunDailyScan(ctx, 0, 0)
}
