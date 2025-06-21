package fetcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/oss-health/background-worker/pkg/db"
	"github.com/oss-health/background-worker/internal/utils"
)

var Resolvers = map[string]func(ctx context.Context, name string) (string, error){
	"npm":  GetNpmRepoURL,
	"pypi": GetPypiRepoURL,
}

func ResolvePendingDependencies(ctx context.Context, batchSize, offset int, ecosystem string) error {
	// requestCtx, cancel := context.WithTimeout(ctxX, 2*time.Minute)
	// defer cancel()

	dependencies, err := db.GetPendingDependencies(ctx, batchSize, offset, ecosystem)
	if err != nil {
		return fmt.Errorf("failed to fetch pending dependencies: %w", err)
	}

	log.Printf("Fetching %d pending dependencies", len(dependencies))

	if len(dependencies) == 0 {
		return nil
	}

	resolvedURLs := make(map[int64]string)
	urlsSet := make(map[string]struct{})
	var resolvedDependencies []db.Dependency
	failureReasons := make(map[int64]string)

	concurrency := 10
	sem := make(chan struct{}, concurrency)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, dependency := range dependencies {
		// NOTE: In Go, the range loop reuses the same dependency variable each time.
		// So if you spin up goroutines in a loop, they all reference the same variable
		// Thus need to make a copy.
		dependencyCopy := dependency
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer func() { 
				<-sem 
				wg.Done()
			}()

			ecosystem := strings.ToLower(dependencyCopy.Ecosystem)
			log.Printf("Resolving url for: %s", dependencyCopy.Name)

			resolver, ok := Resolvers[ecosystem]

			if !ok {
				log.Printf("No resolver found for ecosystem: %s", ecosystem)

				mu.Lock()
				failureReasons[dependencyCopy.ID] = "unsupported ecosystem"
				mu.Unlock()
				return
			}

			if err := utils.WaitUntilAllowed(ctx, ecosystem); err != nil {
				log.Printf("Rate limiter error for %s: %v", dependencyCopy.Name, err)
				return
			}

			url, err := resolver(ctx, dependencyCopy.Name)
			if err != nil {
				log.Printf("Error resolving URL for %s: %v", dependencyCopy.Name, err)

				mu.Lock()
				failureReasons[dependencyCopy.ID] = fmt.Sprintf("resolver error: %v", err)
				mu.Unlock()
				return
			}

			if url == "" {
				log.Printf("No URL found for %s", dependencyCopy.Name)

				mu.Lock()
				failureReasons[dependencyCopy.ID] = "empty URL"
				mu.Unlock()
				return
			}

			mu.Lock()
			resolvedURLs[dependencyCopy.ID] = url
			urlsSet[url] = struct{}{}
			resolvedDependencies = append(resolvedDependencies, dependencyCopy)
			mu.Unlock()
		}()

	}

	wg.Wait()

	if len(resolvedDependencies) == 0 && len(failureReasons) == 0 {
		log.Println("No dependencies resolved or failed")
		return nil
	}

	if len(resolvedDependencies) > 0 {
		urls := make([]string, 0, len(urlsSet))

		for url := range urlsSet {
			urls = append(urls, url)
		}

		urlToID, err := db.UpsertGithubURLs(ctx, urls)
		if err != nil {
			return fmt.Errorf("failed to upsert GitHub URLs: %w", err)
		}

		err = db.BatchUpdateDependencies(ctx, resolvedDependencies, urlToID, resolvedURLs)
		if err != nil {
			return fmt.Errorf("failed to update dependencies: %w", err)
		}

		log.Printf("Processed %d dependencies with GitHub URLs", len(resolvedURLs))
	}

	if len(failureReasons) > 0 {
		err := db.MarkDependenciesAsFailed(ctx, failureReasons)
		if err != nil {
			return fmt.Errorf("failed to mark failed dependencies: %w", err)
		}
		log.Printf("Marked %d dependencies as failed", len(failureReasons))
	}

	return nil
}
