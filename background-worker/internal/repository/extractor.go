package repository

import (
	"log"

	"io/fs"
	"path/filepath"

	"github.com/oss-health/background-worker/internal/repository/parsers"
)

type Extractor interface {
	ExtractDependencies(repositoryPath string) ([]parsers.DependencyParsed, error)
}

type DependencyExtractor struct{}

func (d *DependencyExtractor) ExtractDependencies(repositoryPath string) ([]parsers.DependencyParsed, error) {
	var allDependencies []parsers.DependencyParsed

	err := filepath.WalkDir(repositoryPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
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
					log.Printf("Failed to parse %s: %v\n", path, err)
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
