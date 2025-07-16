package repo

import (
	"context"

	"github.com/oss-health/background-worker/internal/dependencies"
)

type RepoRepositoryInterface interface {
	GetReposForSchedule(ctx context.Context, day int, hour int) ([]Repository, error)
	MarkScanned(ctx context.Context, id int) error
	MarkFailed(ctx context.Context, id int, message string) error
	DeleteRepositoryDependencyVersion(ctx context.Context, repositoryId int) error
	InsertRepositoryDependencyVersion(ctx context.Context, repositoryId int, dependencies []dependencies.DependencyVersion) error
}
