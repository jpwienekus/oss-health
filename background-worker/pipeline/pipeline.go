package pipeline

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/oss-health/background-worker/parser"
	"github.com/oss-health/background-worker/storage"
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
	Store          storage.Store
}

func NewPipeline(fetcher FetcherChain, parserRegistry ParserRegistry, store storage.Store) *Pipeline {
	return &Pipeline{
		FetcherChain:   fetcher,
		ParserRegistry: parserRegistry,
		Store:          store,
	}
}

func (p *Pipeline) Run(repoUrl string) error {
	repositories, err := p.Store.GetRepositoriesBySchedule(0, 0)
	if err != nil {
		return fmt.Errorf("failed to get repositories by schedule: %w", err)
	}

	// TODO: worker pool
	for _, repository := range repositories {
		if err := p.RunRepositoryEnrichment(repository.ID, repository.URL); err != nil {
			log.Printf("failed to run repository enrichment for %s: %v", repository.URL, err)
		}
	}

	return nil
}

func (p *Pipeline) RunRepositoryEnrichment(repositoryId int, repositoryUrl string) error {
	dest, err := osMkdirTemp("", "repo-clone-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	if err := p.cloneRepository(repositoryId, repositoryUrl, dest); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	dependencies, err := p.extractDependencies(dest)
	if err != nil {
		return fmt.Errorf("failed to extract dependencies: %w", err)
	}

	if err := p.insertDependenciesThatDoNotExist(dependencies); err != nil {
		return fmt.Errorf("failed to insert dependencies: %w", err)
	}

	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to remove temporary directory: %w", err)
	}

	return nil
}

func (p *Pipeline) cloneRepository(repositoryId int, repositoryUrl string, destination string) error {
	log.Printf("Cloning repository: %s\n", repositoryUrl)
	if err := p.FetcherChain.Fetch(repositoryUrl, destination); err != nil {
		if err := p.Store.MarkRepositoryAsFailed(repositoryId, err.Error()); err != nil {
			return fmt.Errorf("failed to mark repository as failed: %w", err)
		}

		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	if err := p.Store.MarkRepositoryAsScanned(repositoryId); err != nil {
		return fmt.Errorf("failed to mark repository as scanned: %w", err)
	}

	return nil
}

func (p *Pipeline) extractDependencies(directory string) ([]parser.Dependency, error) {
	log.Printf("Extracting dependencies from directory: %s\n", directory)
	dependencies := []parser.Dependency{}

	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, walkErr error) error {
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

		parsedDependencies := []parser.Dependency{}
		for _, parser := range parsers {
			parsedDeps, err := parser.Parse(content)
			if err != nil {
				return fmt.Errorf("failed to parse file %s with parser %T: %w", path, parser, err)
			}

			parsedDependencies = append(parsedDependencies, parsedDeps...)
		}
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		dependencies = append(dependencies, parsedDependencies...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", directory, err)
	}

	return dependencies, nil
}

func (p *Pipeline) insertDependenciesThatDoNotExist(dependencies []parser.Dependency) error {
	type key struct {
		Name, Ecosystem string
	}

	// These are not specific to a repository, checking what is already stored
	// in the database so no duplicates are inserted
	existingDependencies, err := p.Store.GetExistingDependencies(dependencies)
	if err != nil {
		return fmt.Errorf("failed to get existing dependencies: %w", err)
	}

	existingDependenciesMap := make(map[key]bool, len(existingDependencies))

	for _, dep := range existingDependencies {
		existingDependenciesMap[key{Name: dep.Name, Ecosystem: dep.Ecosystem}] = true
	}

	dependenciesToInsert := make([]parser.Dependency, 0, len(dependencies))

	for _, dep := range dependencies {
		if !existingDependenciesMap[key{Name: dep.Name, Ecosystem: dep.Ecosystem}] {
			dependenciesToInsert = append(dependenciesToInsert, dep)
		}
	}

	if len(dependenciesToInsert) == 0 {
		return nil
	}

	log.Print(dependenciesToInsert)
	if err := p.Store.InsertDependencies(dependenciesToInsert); err != nil {
		return fmt.Errorf("failed to insert dependencies: %w", err)
	}

	return nil
}

func (p *Pipeline) insertVersionsThatDoNotExist(dependencies []parser.Dependency) error {


	return nil
}
