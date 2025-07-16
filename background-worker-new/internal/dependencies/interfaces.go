package dependencies

import (
	"context"

	"github.com/oss-health/background-worker/internal/repo/parsers"
)

type DependencyServiceInterface interface {
	InsertDependenciesThatDoNotExist(ctx context.Context, dependencies []parsers.Dependency, workerId int) error
	InsertVersionsThatDoNotExist(ctx context.Context, dependencies []parsers.Dependency, workerId int) error
	GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]DependencyVersion, error)
}
