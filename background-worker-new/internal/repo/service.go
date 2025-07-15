package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
)

type RepositoryService struct {
	repository RepoRepository
	// dependencyRepository dependency.DependencyRepository
	// cloner               Cloner
	// extractor            Extractor
}

func NewRepositoryService(
	repository RepoRepository,
	// dependencyRepository dependency.DependencyRepository,
	// cloner Cloner,
	// extractor Extractor,
) *RepositoryService {
	return &RepositoryService{
		repository: repository,
		// dependencyRepository: dependencyRepository,
		// cloner:               cloner,
		// extractor:            extractor,
	}
}


func (s *RepositoryService) RunDailyScan(ctx context.Context, day int, hour int) error {
	repositories, err := s.repository.GetReposForSchedule(ctx, day, hour)
	log.Printf("Found %d repositories for schedule: day=%d, hour=%d", len(repositories), day, hour)

	if err != nil {
		return fmt.Errorf("failed to get repositories for schedule: %w", err)
	}

	totalRepositories := len(repositories)
	maxParallelStr := os.Getenv("MAX_PARALLEL")
	maxParallel, err := strconv.Atoi(maxParallelStr)
	if err != nil {
		return fmt.Errorf("invalid MAX_PARALLEL value: %w", err)
	}

	totalParallel := min(maxParallel, totalRepositories)

	workerPool := NewWorkerPool(totalParallel, 2.0, s)
	workerPool.Start(ctx)

	for _, repository := range repositories {
		workerPool.Submit(repository)
	}

	workerPool.Wait()
	log.Printf("Scanning complete")

	return nil
}
