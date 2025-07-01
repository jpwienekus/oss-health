package repository

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/repository/interfaces"
)

type Extractor interface {
	ExtractDependencies(repositoryPath string) ([]dependency.DependencyVersionPair, error)
}

type DependencyExtractor struct {
	Provider interfaces.ParserProvider
}

func (d *DependencyExtractor) ExtractDependencies(repositoryPath string) ([]dependency.DependencyVersionPair, error) {
	depMap := make(map[string]dependency.DependencyVersionPair)

	err := filepath.WalkDir(repositoryPath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk error at %s: %w", path, walkErr)
		}

		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "node_modules", "venv":
				return filepath.SkipDir
			}
			return nil
		}

		parser := d.Provider.GetParser(path)
		if parser == nil {
			return nil
		}

		deps, err := parser.Parse(path)
		if err != nil {
			log.Printf("warning: failed to parse %q: %v", path, err)
			return nil // continue walking
		}

		for _, d := range deps {
			key := fmt.Sprintf("%s@%s", d.Name, d.Version)
			if _, exists := depMap[key]; !exists {
				depMap[key] = d
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("extract dependencies from %q: %w", repositoryPath, err)
	}

	dependencyList := make([]dependency.DependencyVersionPair, 0, len(depMap))

	for _, d := range depMap {
		dependencyList = append(dependencyList, d)
	}

	return dependencyList, nil
}
