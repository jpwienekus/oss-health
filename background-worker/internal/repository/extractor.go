package repository

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/repository/parsers"
)

type Extractor interface {
	ExtractDependencies(repositoryPath string) ([]dependency.DependencyVersionPair, error)
}

type DependencyExtractor struct{}

func (d *DependencyExtractor) ExtractDependencies(repositoryPath string) ([]dependency.DependencyVersionPair, error) {
	depMap := make(map[string]dependency.DependencyVersionPair)
	err := filepath.WalkDir(repositoryPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "node_modules", "venv":
				return filepath.SkipDir
			}
			return nil
		}

		parser := parsers.GetParserForFile(path)

		if parser == nil {
			return nil
		}

		deps, err := parser.Parse(path)

		if err != nil {
			log.Printf("Failed to parse %s: %v", path, err)
			return nil
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
		return nil, err
	}

	allDependencies := make([]dependency.DependencyVersionPair, 0, len(depMap))
	for _, d := range depMap {
		allDependencies = append(allDependencies, d)
	}

	return allDependencies, nil
}
