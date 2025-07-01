package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

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

func (s *RepositoryService) RunDailyScan(ctx context.Context, day int, hour int) error {
	repositories, err := s.repository.GetRepositoriesForDay(ctx, day, hour)

	if err != nil {
		return fmt.Errorf("get repositories for day %d hour %d: %w", day, hour, err)
	}

	totalRepositories := len(repositories)
	maxParallelStr := os.Getenv("MAX_PARALLEL")
	maxParallel, err := strconv.Atoi(maxParallelStr)
	if err != nil {
		return fmt.Errorf("invalid MAX_PARALLEL value %q: %w", maxParallelStr, err)
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

func (s *RepositoryService) CloneAndParse(ctx context.Context, repository Repository) ([]dependency.DependencyVersionPair, error) {
	dependencies, err := s.ProcessRepository(ctx, repository)

	return dependencies, err
}

func (s *RepositoryService) ProcessRepository(ctx context.Context, repo Repository) ([]dependency.DependencyVersionPair, error) {
	tempDir, err := s.cloner.CloneRepository(repo.URL)

	if err != nil {
		if markErr := s.repository.MarkFailed(ctx, repo.ID, err.Error()); markErr != nil {
			return nil, fmt.Errorf("clone %s: %v (and failed to mark as failed: %v)", repo.URL, err, markErr)
		}

		return nil, fmt.Errorf("clone repository %s: %w", repo.URL, err)
	}

	if err := s.repository.MarkScanned(ctx, repo.ID); err != nil {
		return nil, fmt.Errorf("mark repository %d as scanned: %w", repo.ID, err)
	}

	dependencies, err := s.extractor.ExtractDependencies(tempDir)
	if err != nil {
		return nil, fmt.Errorf("extract dependencies from %s: %w", tempDir, err)
	}

	if err := os.RemoveAll(tempDir); err != nil {
		return nil, fmt.Errorf("remove temp dir %s: %w", tempDir, err)
	}

	return dependencies, nil
}

func (s *RepositoryService) ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryId int, pairs []dependency.DependencyVersionPair) error {
	_, err := s.dependencyRepository.ReplaceRepositoryDependencyVersions(ctx, repositoryId, pairs)
	if err != nil {
		return fmt.Errorf("replace dependency versions for repository %d: %w", repositoryId, err)
	}

	return nil
}
