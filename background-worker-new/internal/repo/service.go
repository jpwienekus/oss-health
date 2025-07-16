package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

var ExecCommand = exec.Command
var OsMkdirTemp = os.MkdirTemp

type RepositoryService struct {
	repository        RepoRepositoryInterface
	dependencyService dependencies.DependencyServiceInterface
	ParserRegistry    *parsers.ParserRegistry
}

func NewRepositoryService(
	repository RepoRepositoryInterface,
	dependencyService dependencies.DependencyServiceInterface,
	ParserRegistry *parsers.ParserRegistry,
) *RepositoryService {
	return &RepositoryService{
		repository:        repository,
		dependencyService: dependencyService,
		ParserRegistry:    ParserRegistry,
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

func (s *RepositoryService) CloneRepository(url string, destination string) error {
	cmd := ExecCommand("git", "clone", "--depth=1", url, destination)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v, output: %s", err, output)
	}

	return nil
}

func (s *RepositoryService) ExtractDependencies(ctx context.Context, repo db.Repository) ([]parsers.Dependency, error) {
	tempDir, err := OsMkdirTemp("", "repo-clone-*")
	if err := s.CloneRepository(repo.URL, tempDir); err != nil {
		if err := s.repository.MarkFailed(ctx, repo.ID, fmt.Sprintf("failed to clone repository: %v", err)); err != nil {
			log.Printf("Failed to mark repository %d as failed: %v", repo.ID, err)
		}
		return nil, fmt.Errorf("failed to clone repository %s: %w", repo.URL, err)
	}

	dependencies := []parsers.Dependency{}

	err = filepath.WalkDir(tempDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking the path %q: %w", path, walkErr)
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", ".github", "node_modules", "venv":
				return filepath.SkipDir
			}

			return nil
		}

		base := filepath.Base(path)
		p, ok := s.ParserRegistry.GetParsers(base)

		if !ok {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		parsedDependencies := []parsers.Dependency{}
		for _, parser := range p {
			parsedDeps, err := parser.Parse(content)
			if err != nil {
				return fmt.Errorf("failed to parse file %s with parser %T: %w", path, parser, err)
			}

			parsedDependencies = append(parsedDependencies, parsedDeps...)
		}

		dependencies = append(dependencies, parsedDependencies...)

		return nil
	})

	if err != nil {
		if err := s.repository.MarkFailed(ctx, repo.ID, fmt.Sprintf("failed to walk directory: %v", err)); err != nil {
			return nil, fmt.Errorf("failed to mark repository %d as failed: %w", repo.ID, err)
		}

		return nil, fmt.Errorf("failed to walk directory %s: %w", tempDir, err)
	}

	if err := os.RemoveAll(tempDir); err != nil {
		return nil, fmt.Errorf("failed to remove temporary directory %s: %w", tempDir, err)
	}

	if err := s.repository.MarkScanned(ctx, repo.ID); err != nil {
		return nil, fmt.Errorf("failed to mark repository %d as scanned: %w", repo.ID, err)
	}

	return dependencies, nil
}

func (s *RepositoryService) ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryID int, pairs []parsers.Dependency, workerId int) error {
	if err := s.dependencyService.InsertDependenciesThatDoNotExist(ctx, pairs, workerId); err != nil {
		return fmt.Errorf("failed to insert dependencies that do not exist: %w", err)
	}

	if err := s.dependencyService.InsertVersionsThatDoNotExist(ctx, pairs, workerId); err != nil {
		return fmt.Errorf("failed to insert versions that do not exist: %w", err)
	}

	if err := s.saveRepositoryDependencyVersions(ctx, repositoryID, pairs); err != nil {
		return fmt.Errorf("failed to save repository dependency versions: %w", err)
	}

	return nil
}

func (s *RepositoryService) saveRepositoryDependencyVersions(ctx context.Context, repositoryId int, dependencies []parsers.Dependency) error {
	if err := s.repository.DeleteRepositoryDependencyVersion(ctx, repositoryId); err != nil {
		return fmt.Errorf("failed to delete repository dependency versions: %w", err)
	}

	existingVersions, err := s.dependencyService.GetExistingVersionsForDependencies(ctx, dependencies)
	if err != nil {
		return fmt.Errorf("failed to get existing versions: %w", err)
	}

	if len(existingVersions) == 0 {
		return nil
	}

	if err := s.repository.InsertRepositoryDependencyVersion(ctx, repositoryId, existingVersions); err != nil {
		return fmt.Errorf("failed to insert repository dependency versions: %w", err)
	}

	return nil
}
