package fetcher

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/oss-health/background-worker/internal/utils"
	"github.com/oss-health/background-worker/internal/db"
)

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}
var Resolvers = map[string]func(ctx context.Context, name string) (string, error){
	"npm":  GetNpmRepoURL(defaultHTTPClient, "https://registry.npmjs.org"),
	"pypi": GetPypiRepoURL(defaultHTTPClient, "https://pypi.org/pypi"),
}

func ResolvePendingDependencies(
	ctx context.Context,
	batchSize int,
	offset int,
	ecosystem string,
	rateLimiter utils.RateLimiter,
	resolvers map[string]func(ctx context.Context, name string) (string, error),
) error {
	dependencies, err := db.GetPendingDependencies(ctx, batchSize, offset, ecosystem)

	if err != nil {
		return fmt.Errorf("failed to fetch pending dependencies: %w", err)
	}

	log.Printf("Fetching %d pending dependencies", len(dependencies))

	if len(dependencies) == 0 {
		return nil
	}

	var resolvedDependencies []db.Dependency
	resolvedURLs := make(map[int64]string)
	urlsSet := make(map[string]struct{})
	failureReasons := make(map[int64]string)

	concurrency := 10
	semaphore := make(chan struct{}, concurrency)
	var mutex sync.Mutex
	var waitGroup sync.WaitGroup

	for _, dependency := range dependencies {
		// NOTE: In Go, the range loop reuses the same dependency variable each time.
		// So if you spin up goroutines in a loop, they all reference the same variable
		// Thus need to make a copy.
		dependencyCopy := dependency
		semaphore <- struct{}{}
		waitGroup.Add(1)

		go func() {
			defer func() {
				<-semaphore
				waitGroup.Done()
			}()

			ecosystem := strings.ToLower(dependencyCopy.Ecosystem)
			log.Printf("Resolving url for: %s", dependencyCopy.Name)

			resolver, ok := resolvers[ecosystem]

			if !ok {
				log.Printf("No resolver found for ecosystem: %s", ecosystem)

				mutex.Lock()
				failureReasons[dependencyCopy.ID] = "unsupported ecosystem"
				mutex.Unlock()
				return
			}

			if err := rateLimiter.WaitUntilAllowed(ctx, ecosystem); err != nil {
				log.Printf("Rate limiter error for %s: %v", dependencyCopy.Name, err)
				return
			}

			url, err := resolver(ctx, dependencyCopy.Name)
			if err != nil {
				log.Printf("Error resolving URL for %s: %v", dependencyCopy.Name, err)

				mutex.Lock()
				failureReasons[dependencyCopy.ID] = fmt.Sprintf("resolver error: %v", err)
				mutex.Unlock()
				return
			}

			if url == "" {
				log.Printf("No URL found for %s", dependencyCopy.Name)

				mutex.Lock()
				failureReasons[dependencyCopy.ID] = "empty URL"
				mutex.Unlock()
				return
			}

			mutex.Lock()
			resolvedURLs[dependencyCopy.ID] = url
			urlsSet[url] = struct{}{}
			resolvedDependencies = append(resolvedDependencies, dependencyCopy)
			mutex.Unlock()
		}()

	}

	waitGroup.Wait()

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
