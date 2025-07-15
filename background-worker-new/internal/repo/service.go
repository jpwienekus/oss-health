package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/oss-health/background-worker/internal/repo/parser"
)

var execCommand = exec.Command
var osMkdirTemp = os.MkdirTemp

type RepositoryService struct {
	repository RepoRepository

	ParserRegistry *parser.ParserRegistry
	// dependencyRepository dependency.DependencyRepository
	// cloner               Cloner
	// extractor            Extractor
}

func NewRepositoryService(
	repository RepoRepository,
	ParserRegistry *parser.ParserRegistry,
) *RepositoryService {
	return &RepositoryService{
		repository:     repository,
		ParserRegistry: ParserRegistry,
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
	cmd := execCommand("git", "clone", "--depth=1", url, destination)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v, output: %s", err, output)
	}

	return nil
}

func (s *RepositoryService) ExtractDependencies(ctx context.Context, url string) ([]parser.Dependency, error) {
	tempDir, err := osMkdirTemp("", "repo-clone-*")
	if err := s.CloneRepository(url, tempDir); err != nil {
		return nil, fmt.Errorf("failed to clone repository %s: %w", url, err)
	}

	log.Printf("Extracting dependencies from directory: %s\n", tempDir)
	dependencies := []parser.Dependency{}

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
		parsers, ok := s.ParserRegistry.GetParsers(base)

		if !ok {
			return nil
		}

		log.Printf("Processing file: %s\n", path)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		parsedDependencies := []parser.Dependency{}
		for _, parser := range parsers {
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
		return nil, fmt.Errorf("failed to walk directory %s: %w", tempDir, err)
	}

	if err := os.RemoveAll(tempDir); err != nil {
		return nil, fmt.Errorf("failed to remove temporary directory %s: %w", tempDir, err)
	}

	return dependencies, nil
}
