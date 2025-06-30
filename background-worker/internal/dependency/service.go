package dependency

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	resolvedUrlMap := make(map[int64]string)
	failureReasons := make(map[int64]string)

	for _, dependency := range dependencies {
		ecosystem := strings.ToLower(dependency.Ecosystem)
		resolver, ok := s.resolvers[ecosystem]

		if !ok || resolver == nil {
			failureReasons[dependency.ID] = "unsupported ecosystem"
			continue
		}

		if err := s.rateLimiter.WaitUntilAllowed(ctx, ecosystem); err != nil {
			failureReasons[dependency.ID] = fmt.Sprintf("rate limiter error: %v", err)
			continue
		}

		url, err := resolver(ctx, dependency.Name)

		if err != nil {
			failureReasons[dependency.ID] = fmt.Sprintf("resolver error: %v", err)
			continue
		}

		if url == "" {
			failureReasons[dependency.ID] = "empty URL"
			continue
		}

		resolvedUrlMap[dependency.ID] = url
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
