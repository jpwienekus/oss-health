package pipeline

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/oss-health/background-worker/parser"
)

var osMkdirTemp = os.MkdirTemp

type FetcherChain interface {
	Fetch(repoUrl string, destPath string) error
}

type ParserRegistry interface {
	GetParsers(filename string) ([]parser.Parser, bool)
}

type Pipeline struct {
	FetcherChain   FetcherChain
	ParserRegistry ParserRegistry
}

func NewPipeline(fetcher FetcherChain, parserRegistry ParserRegistry) *Pipeline {
	return &Pipeline{
		FetcherChain:   fetcher,
		ParserRegistry: parserRegistry,
	}
}

func (p *Pipeline) Run(repoUrl string) error {
	dest, err := osMkdirTemp("", "repo-clone-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	log.Printf("Cloning repository: %s\n", repoUrl)
	if err := p.FetcherChain.Fetch(repoUrl, dest); err != nil {
		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	err = filepath.WalkDir(dest, func(path string, d os.DirEntry, walkErr error) error {
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
		parsers, ok := p.ParserRegistry.GetParsers(base)

		if !ok {
			return nil
		}

		log.Printf("Processing file: %s\n", path)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		dependencies := []parser.Dependency{}
		for _, parser := range parsers {
			parsedDeps, err := parser.Parse(content)
			if err != nil {
				return fmt.Errorf("failed to parse file %s with parser %T: %w", path, parser, err)
			}

			dependencies = append(dependencies, parsedDeps...)
		}
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		log.Print(dependencies)

		return nil
	})

	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to remove temporary directory: %w", err)
	}

	return nil
}
