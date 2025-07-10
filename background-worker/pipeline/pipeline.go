package pipeline

import (
	"fmt"
	"log"
	"os"
)

var osMkdirTemp = os.MkdirTemp

type FetcherChain interface {
	Fetch(repoUrl string, destPath string) error
}

type Pipeline struct {
	FetcherChain FetcherChain
}

func NewPipeline(fetcher FetcherChain) *Pipeline {
	return &Pipeline{
		FetcherChain: fetcher,
	}
}

func (p *Pipeline) Run(repoUrl string) error {
	dest, err := osMkdirTemp("", "repo-clone-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	if err := p.FetcherChain.Fetch(repoUrl, dest); err != nil {
		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	log.Printf("Repository cloned to: %s\n", dest)

	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to remove temporary directory: %w", err)
	}

	return nil
}
