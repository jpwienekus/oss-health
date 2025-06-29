package repository

import (
	"context"
	// "fmt"
	"log"
	"os"

	"github.com/oss-health/background-worker/internal/dependency"
)

type RepositoryService struct {
	repository           RepositoryRepository
	dependencyRepository dependency.DependencyRepository
	cloner               Cloner
	extractor            Extractor
}

func NewRepositoryService(
	repository RepositoryRepository,
	dependencyRepository dependency.DependencyRepository,
	cloner Cloner,
	extractor Extractor,
) *RepositoryService {
	return &RepositoryService{
		repository:           repository,
		dependencyRepository: dependencyRepository,
		cloner:               cloner,
		extractor:            extractor,
	}
}

func (s *RepositoryService) RunDailyScan(ctx context.Context, day int, hour int) {
	repositories, err := s.repository.GetRepositoriesForDay(ctx, day, hour)

	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
	}

	totalRepositories := len(repositories)
	// TODO: make env var
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

func (s *RepositoryService) CloneAndParse(ctx context.Context, repository Repository) ([]dependency.DependencyVersionPair, error) {
	dependencies, err := s.ProcessRepository(ctx, repository)

	return dependencies, err
}

func (s *RepositoryService) ProcessRepository(ctx context.Context, repo Repository) ([]dependency.DependencyVersionPair, error) {
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

func (s *RepositoryService) ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryId int, pairs []dependency.DependencyVersionPair) {
	_, err := s.dependencyRepository.ReplaceRepositoryDependencyVersions(ctx, repositoryId, pairs)

	if err != nil {
		log.Print(err)
	}
}
