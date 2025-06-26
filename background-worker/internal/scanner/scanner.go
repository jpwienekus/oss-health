package scanner

import (
	"context"
	"log"

	"github.com/oss-health/background-worker/internal/repository"
)

func RunDailyScan(ctx context.Context, day int, hour int) {
	repositories, err := repository.GetRepositoriesForDay(ctx, day, hour)

	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
	}

	log.Printf("Scanning %d repositories for day %d hour %d", len(repositories), day, hour)
	workerPool := NewWorkerPool(16, 2.0)
	workerPool.Start(ctx)

	for _, repository := range repositories {
		workerPool.Submit(repository)
	}

	workerPool.Wait()
	log.Printf("Scanning complete")
}
