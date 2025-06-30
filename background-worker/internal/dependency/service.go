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
	// dependency    Dependency
}

func (s *DependencyService) ResolvePendingDependencies(ctx context.Context, batchSize int, offset int, ecosystem string) error {
	// TODO: include a date stale check here as part of GetDependenciesPendingUrlResolution
	// The upsert logic below supports handling existing entries
	dependencies, err := s.repository.GetDependenciesPendingUrlResolution(ctx, batchSize, offset, ecosystem)

	if err != nil {
		return fmt.Errorf("failed to fetch pending dependencies: %w", err)
	}

	log.Printf("Resolving urls for %d dependencies (%s)", len(dependencies), ecosystem)

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
			return ctx.Err()
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

			resultsCh <- resolveResult{id: dependencyCopy.ID, url: url }
		}()
	}

	wg.Wait()
	close(resultsCh)

	if len(resultsCh) == 0 {
		return nil
	}

	resolvedUrlMap := make(map[int64]string)
	failureReasons := make(map[int64]string)

	for result := range resultsCh {
		if result.failureReason != "" {
			failureReasons[result.id] = result.failureReason
		} else {
			resolvedUrlMap[result.id] = result.url
		}
	}

	if len(resolvedUrlMap) > 0 {
		dependencyDependencyRepositoryIdMap, err := s.repository.UpsertGithubURLs(ctx, resolvedUrlMap)

		if err != nil {
			return fmt.Errorf("failed to upsert GitHub URLs: %w", err)
		}

		err = s.repository.BatchUpdateDependencies(ctx, dependencyDependencyRepositoryIdMap)

		if err != nil {
			return fmt.Errorf("failed to update dependency with dependency repository link: %w", err)
		}

		log.Printf("Processed %d dependencies with GitHub URLs", len(resolvedUrlMap))
	}

	if len(failureReasons) > 0 {
		err := s.repository.MarkDependenciesAsFailed(ctx, failureReasons)

		if err != nil {
			return fmt.Errorf("failed to mark failed dependencies: %w", err)
		}

		log.Printf("Marked %d dependencies as failed", len(failureReasons))
	}

	return nil
}
