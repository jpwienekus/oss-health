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

	for _, dependency := range dependencies {
		ecosystem := strings.ToLower(dependency.Ecosystem)

		resolver, ok := Resolvers[ecosystem]
		if !ok {
			log.Printf("No resolver found for ecosystem: %s", ecosystem)
			continue
		}

		if err := utils.WaitUntilAllowed(ctx, ecosystem); err != nil {
			log.Printf("Rate limiter error for %s: %v", dependency.Name, err)
			continue
		}

		url, err := resolver(ctx, dependency.Name)
		if err != nil {
			log.Printf("Error resolving URL for %s: %v", dependency.Name, err)
			continue
		}

		if url == "" {
			log.Printf("No URL found for %s", dependency.Name)
			continue
		}

		if err := db.UpdateDependencyURL(ctx, dependency.ID, url); err != nil {
			log.Printf("Error updating DB for %s: %v", dependency.Name, err)
		} else {
			log.Printf("Updated dependency %s with URL: %s", dependency.Name, url)
		}
	}

	return nil
}
