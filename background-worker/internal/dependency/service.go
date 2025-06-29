package dependency

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/oss-health/background-worker/internal/utils"
)

type DependencyService struct {
	repository DependencyRepository

	rateLimiter utils.RateLimiter
	resolvers   map[string]func(ctx context.Context, name string) (string, error)
}

func NewDependencyService(
	repository DependencyRepository,
	limiters utils.RateLimiter,
	resolvers map[string]func(ctx context.Context, name string) (string, error),
) *DependencyService {
	return &DependencyService{
		repository:  repository,
		rateLimiter: limiters,
		resolvers:   resolvers,
	}
}

type resolveResult struct {
	id            int64
	url           string
	failureReason string
	dependency    Dependency
}

func (s *DependencyService) ResolvePendingDependencies(ctx context.Context, batchSize int, offset int, ecosystem string) error {
	dependencies, err := s.repository.GetPendingDependencies(ctx, batchSize, offset, ecosystem)

	if err != nil {
		return fmt.Errorf("failed to fetch pending dependencies: %w", err)
	}

	log.Printf("Resolving urls for %d dependencies", len(dependencies))

	if len(dependencies) == 0 {
		return nil
	}

	concurrency := 10
	semaphore := make(chan struct{}, concurrency)
	resultsCh := make(chan resolveResult, len(dependencies))

	var wg sync.WaitGroup

	for _, dependency := range dependencies {
		select {
		case <-ctx.Done():
			break
		case semaphore <- struct{}{}:
		}

		wg.Add(1)
		// NOTE: In Go, the range loop reuses the same dependency variable each time.
		// So if you spin up goroutines in a loop, they all reference the same variable
		// Thus need to make a copy.
		dependencyCopy := dependency

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()

			ecosystem := strings.ToLower(dependencyCopy.Ecosystem)
			resolver, ok := s.resolvers[ecosystem]

			if !ok || resolver == nil {
				resultsCh <- resolveResult{id: dependencyCopy.ID, failureReason: "unsupported ecosystem"}
				return
			}

			if err := s.rateLimiter.WaitUntilAllowed(ctx, ecosystem); err != nil {
				resultsCh <- resolveResult{id: dependencyCopy.ID, failureReason: fmt.Sprintf("rate limiter error: %v", err)}
				return
			}

			url, err := resolver(ctx, dependencyCopy.Name)

			if err != nil {
				resultsCh <- resolveResult{id: dependencyCopy.ID, failureReason: fmt.Sprintf("resolver error: %v", err)}
				return
			}
			if url == "" {
				resultsCh <- resolveResult{id: dependencyCopy.ID, failureReason: "empty URL"}
				return
			}

			resultsCh <- resolveResult{id: dependencyCopy.ID, url: url, dependency: dependencyCopy}
		}()
	}

	wg.Wait()
	close(resultsCh)

	resolvedDeps := make([]Dependency, 0)
	resolvedURLs := make(map[int64]string)
	urlSet := make(map[string]struct{})
	failureReasons := make(map[int64]string)

	for res := range resultsCh {
		if res.failureReason != "" {
			failureReasons[res.id] = res.failureReason
		} else {
			resolvedURLs[res.id] = res.url
			urlSet[res.url] = struct{}{}
			resolvedDeps = append(resolvedDeps, res.dependency)
		}
	}

	if len(resolvedDeps) > 0 {
		urls := make([]string, 0, len(urlSet))
		for url := range urlSet {
			urls = append(urls, url)
		}

		urlToID, err := s.repository.UpsertGithubURLs(ctx, urls)
		if err != nil {
			return fmt.Errorf("failed to upsert GitHub URLs: %w", err)
		}

		if err := s.repository.BatchUpdateDependencies(ctx, resolvedDeps, urlToID, resolvedURLs); err != nil {
			return fmt.Errorf("failed to update dependencies: %w", err)
		}
		log.Printf("Processed %d dependencies with GitHub URLs", len(resolvedURLs))
	}

	if len(failureReasons) > 0 {
		if err := s.repository.MarkDependenciesAsFailed(ctx, failureReasons); err != nil {
			return fmt.Errorf("failed to mark failed dependencies: %w", err)
		}
		log.Printf("Marked %d dependencies as failed", len(failureReasons))
	}

	return nil
}
