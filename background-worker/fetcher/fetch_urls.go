package fetcher

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/oss-health/background-worker/db"
	"github.com/oss-health/background-worker/utils"
)

var Resolvers = map[string]func(ctx context.Context, name string) (string, error){
	"npm":  GetNpmRepoURL,
	"pypi": GetPypiRepoURL,
}

func ResolvePendingDependencies(ctx context.Context, batchSize, offset int, ecosystem string) error {
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

	for _, dependency := range dependencies {
		ecosystem := strings.ToLower(dependency.Ecosystem)
		log.Printf("Resolving url for: %s", dependency.Name)

		resolver, ok := Resolvers[ecosystem]
		if !ok {
			log.Printf("No resolver found for ecosystem: %s", ecosystem)
			failureReasons[dependency.ID] = "unsupported ecosystem"
			continue
		}

		if err := utils.WaitUntilAllowed(ctx, ecosystem); err != nil {
			log.Printf("Rate limiter error for %s: %v", dependency.Name, err)
			continue
		}

		url, err := resolver(ctx, dependency.Name)
		if err != nil {
			log.Printf("Error resolving URL for %s: %v", dependency.Name, err)
			failureReasons[dependency.ID] = fmt.Sprintf("resolver error: %v", err)
			continue
		}

		if url == "" {
			log.Printf("No URL found for %s", dependency.Name)
			failureReasons[dependency.ID] = "empty URL"
			continue
		}

		resolvedURLs[dependency.ID] = url
		urlsSet[url] = struct{}{}
		resolvedDependencies = append(resolvedDependencies, dependency)
	}

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
