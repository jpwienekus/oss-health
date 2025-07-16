package dependencies

import (
	"context"

	"github.com/oss-health/background-worker/internal/repo/parsers"
)

type DependencyRepository interface {
	GetExistingDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]Dependency, error)
	GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]DependencyVersion, error)
	InsertDependencies(ctx context.Context, dependencies []parsers.Dependency) error
	InsertDependencyVersions(ctx context.Context, dependencies []DependencyVersion) error
}
