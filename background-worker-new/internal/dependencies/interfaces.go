package dependencies

import (
	"context"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

type DependencyServiceInterface interface {
	InsertDependenciesThatDoNotExist(ctx context.Context, dependencies []parsers.Dependency, workerId int) error
	InsertVersionsThatDoNotExist(ctx context.Context, dependencies []parsers.Dependency, workerId int) error
	GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.DependencyVersion, error)
	ResolveDependencyUrls(ctx context.Context, repositoryId int) error
}

type DependencyRepositoryInterface interface {
	GetExistingDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.Dependency, error)
	GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.DependencyVersion, error)
	InsertDependencies(ctx context.Context, dependencies []parsers.Dependency) error
	InsertDependencyVersions(ctx context.Context, dependencies []db.DependencyVersion) error
	GetDependenciesPendingUrlResolution(ctx context.Context, repositoryId int) ([]db.Dependency, error)
	InsertDependencyRepository(ctx context.Context, dependencyUrlMap map[int]string) error
	GetDependencyRepositoryIds(ctx context.Context, dependencyUrlMap map[int]string) (map[string]int, error)
	UpdateDependencyRepositoryIDs(ctx context.Context, depRepoMap map[int]int) error
	MarkDependenciesAsFailed(ctx context.Context, failureReasons map[int]string) error
	GetDependenciesPendingOSVScan(ctx context.Context) ([]db.DependencyVersion, error)
}

type ResolverChainInterface interface {
	Resolve(ctx context.Context, dependency db.Dependency) (string, error)
}
