package storage

import (
	"github.com/oss-health/background-worker/internal/models"
	"github.com/oss-health/background-worker/parser"
)

type Store interface {
	GetRepositoriesBySchedule(day int, hour int) ([]models.Repository, error)
	MarkRepositoryAsScanned(id int) error
	MarkRepositoryAsFailed(id int, message string) error
	DeleteRepositoryDependencyVersion(repositoryId int) error
	InsertRepositoryDependencyVersion(repositoryId int, dependencyVersion []models.DependencyVersion) error

	GetExistingDependencies(dependencies []parser.Dependency) ([]models.Dependency, error)
	GetDependenciesPendingUrlResolution(repositoryId int) ([]models.Dependency, error)
	InsertDependencies(dependencies []parser.Dependency) error
	UpdateDependencyRepositoryIDs(dependencyDependencyRepository map[int]int) error
	MarkDependenciesAsFailed(failureReasons map[int]string) error

	GetExistingVersionsForDependencies(dependencies []parser.Dependency) ([]models.DependencyVersion, error)
	InsertDependencyVersions(dependencies []models.DependencyVersion) error

	GetDependencyRepositoryIds(dependencyUrlMap map[int]string) (map[string]int, error)
	InsertDependencyRepository(dependencyUrlMap map[int]string) error
}
