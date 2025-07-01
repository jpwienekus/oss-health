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
		dependencyDependencyRepositoryIdMap, err := s.repository.UpsertGithubURLs(ctx, resolvedUrls)

		if err != nil {
			return fmt.Errorf("upsert GitHub URLs: %w", err)
		}

		if err := s.repository.BatchUpdateDependencies(ctx, dependencyDependencyRepositoryIdMap); err != nil {
			return fmt.Errorf("update dependencies: %w", err)
		}

		log.Printf("Processed %d dependencies with GitHub URLs", len(resolvedUrls))
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
	resolved := make(map[int64]string)
	failures := make(map[int64]string)

	for _, dep := range dependencies {
		ecosystem := strings.ToLower(dep.Ecosystem)
		resolver, ok := s.resolvers[ecosystem]

		if !ok || resolver == nil {
			failures[dep.ID] = "unsupported ecosystem"
			continue
		}

		if err := s.rateLimiter.WaitUntilAllowed(ctx, ecosystem); err != nil {
			failures[dep.ID] = fmt.Sprintf("rate limit: %v", err)
			continue
		}

		url, err := resolver(ctx, dep.Name)
		switch {
		case err != nil:
			failures[dep.ID] = fmt.Sprintf("resolve: %v", err)
		case url == "":
			failures[dep.ID] = "empty URL"
		default:
			resolved[dep.ID] = url
		}
	}

	return resolved, failures
}
