package dependency

import (
	"context"
)

type DependencyRepository interface {
	GetDependenciesPendingUrlResolution(ctx context.Context, batchSize, offset int, ecosystem string) ([]Dependency, error)
	UpsertGithubURLs(ctx context.Context, urls []string) (map[string]int64, error)
	BatchUpdateDependencies(ctx context.Context, dependencies []Dependency, urlToID map[string]int64, resolvedURLs map[int64]string) error
	MarkDependenciesAsFailed(ctx context.Context, failureReasons map[int64]string) error
	ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryID int, pairs []DependencyVersionPair) ([]DependencyVersionResult, error)
}
