package dependency

import (
	"context"
)

type DependencyRepository interface {
	GetDependenciesPendingUrlResolution(ctx context.Context, batchSize, offset int, ecosystem string) ([]Dependency, error)
	UpsertGithubURLs(ctx context.Context, resolvedUrls map[int64]string) (map[int64]int64, error)
	BatchUpdateDependencies(ctx context.Context, dependencyDependencyRepositoryIdMap map[int64]int64) error
	MarkDependenciesAsFailed(ctx context.Context, failureReasons map[int64]string) error
	ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryID int, pairs []DependencyVersionPair) ([]DependencyVersionResult, error)
}
