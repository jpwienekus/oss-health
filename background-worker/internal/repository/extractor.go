package repository

import (
	"log"

	"io/fs"
	"path/filepath"
)

type Extractor interface {
	ExtractDependencies(repositoryPath string) ([]DependencyParsed, error)
}

type DependencyExtractor struct{}

func (d *DependencyExtractor) ExtractDependencies(repositoryPath string) ([]DependencyParsed, error) {
	var allDependencies []DependencyParsed

	err := filepath.WalkDir(repositoryPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}

		filename := filepath.Base(path)

		for _, parser := range DependencyParsers {
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
