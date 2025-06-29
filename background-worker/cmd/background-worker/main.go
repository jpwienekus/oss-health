package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repository"
)

func main() {
	// TODO: env var
	connectionString := "postgres://dev-user:password@localhost:5432/dev_db"

	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepositoryRepository(db)
	cloner := &repository.GitCloner{}
	extractor := &repository.DependencyExtractor{}
	service := repository.NewRepositoryService(repo, cloner, extractor)
	// service.RunDailyScan(ctx, 4, 3)
	service.RunDailyScan(ctx, 6, 9)

	// scheduler.Start()

	// Graceful shutdown block
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down gracefully.")
}
