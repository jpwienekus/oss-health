package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oss-health/background-worker/fetcher"
)

func main() {
	repoUrl := "https://github.com/jpwienekus/oss-health.git"
	dest, err := os.MkdirTemp("", "repo-clone-*")
	if err != nil {
		log.Fatalf("failed to create temp directory: %v", err)
	}

	fetcherChain := fetcher.NewFetcherChain()
	if err := fetcherChain.Fetch(repoUrl, dest); err != nil {
		log.Fatalf("failed to fetch repository: %v", err)
	}

	fmt.Printf("Repository cloned to: %s\n", dest)
}
