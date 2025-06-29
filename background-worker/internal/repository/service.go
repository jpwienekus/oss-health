package repository

import (
	"context"
	"log"
	"os"

	"github.com/oss-health/background-worker/internal/repository/parsers"
)

type RepositoryService struct {
	repository RepositoryRepository
	cloner     Cloner
	extractor  Extractor
}

func NewRepositoryService(
	repository RepositoryRepository,
	cloner Cloner,
	extractor Extractor,
) *RepositoryService {
	return &RepositoryService{
		repository: repository,
		cloner:     cloner,
		extractor:  extractor,
	}
}

func (s *RepositoryService) RunDailyScan(ctx context.Context, day int, hour int) {
	repositories, err := s.repository.GetRepositoriesForDay(ctx, day, hour)

	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
	}

	totalRepositories := len(repositories)
	log.Printf("Scanning %d repositories for day %d hour %d", totalRepositories, day, hour)
	maxParallel := 16
	totalParallel := min(maxParallel, totalRepositories)

	workerPool := NewWorkerPool(totalParallel, 2.0, s)
	workerPool.Start(ctx)

	for _, repository := range repositories {
		workerPool.Submit(repository)
	}

	workerPool.Wait()
	log.Printf("Scanning complete")
}

func (s *RepositoryService) CloneAndParse(ctx context.Context, repository Repository) {
	dependencies, err := s.ProcessRepository(ctx, repository)

	if err != nil {
		log.Printf("Could not parse dependencies for %s: %v", repository.URL, err)
	}

	print("Dependencies:")
	print(dependencies)
}

func (s *RepositoryService) ProcessRepository(ctx context.Context, repo Repository) ([]parsers.DependencyParsed, error) {
	tempDir, err := s.cloner.CloneRepository(repo.URL)

	if err != nil {
		log.Printf("Failed to process %s: %v", repo.URL, err)
		s.repository.MarkFailed(ctx, repo.ID, err.Error())
		return nil, err
	} else {
		s.repository.MarkScanned(ctx, repo.ID)
	}

	dependencies, err := s.extractor.ExtractDependencies(tempDir)

	if err != nil {
		return nil, err
	}

	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("failed to remove %s: %v", tempDir, err)
		return nil, err
	}

	return dependencies, nil
}
