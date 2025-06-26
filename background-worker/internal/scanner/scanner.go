package scanner

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"io/fs"
	"path/filepath"

	"github.com/oss-health/background-worker/internal/parsers"
	"github.com/oss-health/background-worker/internal/repository"
)

func RunDailyScan(ctx context.Context, day int, hour int) {
	repositories, err := repository.GetRepositoriesForDay(ctx, day, hour)

	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
	}

	for _, repository := range repositories {
		go processRepository(ctx, repository)
	}
}

func processRepository(ctx context.Context, repo repository.Repository) {
	tempDir, err := cloneRepository(repo.URL)

	if err != nil {
		log.Printf("Failed to process %s: %v", repo.URL, err)
		repository.MarkFailed(ctx, repo.ID, err.Error())
	} else {
		repository.MarkScanned(ctx, repo.ID)
	}

	dependencies, err := extractDependencies(tempDir)

	if err != nil {
		log.Printf("Could not get dependencies for %s: %v", repo.URL, err)
	}
	log.Print(len(dependencies))

	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("failed to remove %s: %v", tempDir, err)
	}
}

func cloneRepository(url string) (string, error) {
	tempDir, err := os.MkdirTemp("", "repo-clone-*")

	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth", "1", url, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if err := os.RemoveAll(tempDir); err != nil {
			return "", fmt.Errorf("failed to remove %s: %w", tempDir, err)
		}

		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return tempDir, nil

}

func extractDependencies(repositoryPath string) ([]parsers.Dependency, error) {
	var allDependencies []parsers.Dependency

	err := filepath.WalkDir(repositoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		filename := filepath.Base(path)

		for _, parser := range parsers.DependencyParsers {
			match, err := filepath.Match(parser.Pattern, filename)
			if err != nil {
				continue
			}
			if !match {
				match, _ = filepath.Match(parser.Pattern, path)
			}
			if match {
				deps, err := parser.Parse(path)
				if err != nil {
					fmt.Printf("Failed to parse %s: %v\n", path, err)
					continue
				}
				allDependencies = append(allDependencies, deps...)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return allDependencies, nil
}
