package repo

import (
	"context"

	"github.com/oss-health/background-worker/internal/db"
)

type RepoRepositoryInterface interface {
	GetReposForSchedule(ctx context.Context, day int, hour int) ([]db.Repository, error)
	MarkScanned(ctx context.Context, id int) error
	MarkFailed(ctx context.Context, id int, message string) error
	DeleteRepositoryDependencyVersion(ctx context.Context, repositoryId int) error
	InsertRepositoryDependencyVersion(ctx context.Context, repositoryId int, dependencies []db.DependencyVersion) error
}
