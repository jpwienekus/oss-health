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

func (s *DependencyService) ResolvePendingDependencies(ctx context.Context, batchSize int, offset int, ecosystem string) error {
	// TODO: include a date stale check here as part of GetDependenciesPendingUrlResolution
	// The upsert logic below supports handling existing entries
	dependencies, err := s.repository.GetDependenciesPendingUrlResolution(ctx, batchSize, offset, ecosystem)
	if err != nil {
		return fmt.Errorf("get pending dependencies: %w", err)
	}

	log.Printf("Resolving urls for %d dependencies (%s)", len(dependencies), ecosystem)

	if len(dependencies) == 0 {
		return nil
	}

	resolvedUrls, failures := s.resolveURLs(ctx, dependencies)

	if len(resolvedUrls) > 0 {
		dependencyDependencyRepositoryIdMap, err := s.repository.UpsertRepositoryURLs(ctx, resolvedUrls)

		if err != nil {
			return fmt.Errorf("upsert Repository URLs: %w", err)
		}

		if err := s.repository.BatchUpdateDependencies(ctx, dependencyDependencyRepositoryIdMap); err != nil {
			return fmt.Errorf("update dependencies: %w", err)
		}

		log.Printf("Processed %d dependencies with Repository URLs", len(resolvedUrls))
	}

	if len(failures) > 0 {
		err := s.repository.MarkDependenciesAsFailed(ctx, failures)

		if err != nil {
			return fmt.Errorf("mark failed dependencies: %w", err)
		}

		log.Printf("Marked %d dependencies as failed", len(failures))
	}

	return nil
}

func (s *DependencyService) resolveURLs(ctx context.Context, dependencies []Dependency) (map[int64]string, map[int64]string) {
	var (
		resolved = make(map[int64]string)
		failures = make(map[int64]string)
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	for _, dep := range dependencies {
		depCopy := dep
		wg.Add(1)

		go func() {
			defer wg.Done()

			ecosystem := strings.ToLower(depCopy.Ecosystem)
			resolver, ok := s.resolvers[ecosystem]

			if !ok || resolver == nil {
				mu.Lock()
				failures[depCopy.ID] = "unsupported ecosystem"
				mu.Unlock()
				return
			}

			if err := s.rateLimiter.WaitUntilAllowed(ctx, ecosystem); err != nil {
				mu.Lock()
				failures[depCopy.ID] = fmt.Sprintf("rate limit: %v", err)
				mu.Unlock()
				return
			}

			url, err := resolver(ctx, depCopy.Name)
			log.Print(url)

			mu.Lock()
			defer mu.Unlock()
			switch {
			case err != nil:
				failures[depCopy.ID] = fmt.Sprintf("resolve: %v", err)
			case url == "":
				failures[depCopy.ID] = "empty URL"
			default:
				resolved[depCopy.ID] = url
			}
		}()
	}

	wg.Wait()

	return resolved, failures
}
