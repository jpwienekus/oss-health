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

	GetExistingDependencies(dependencies []parser.Dependency) ([]models.Dependency, error)
	InsertDependencies(dependencies []parser.Dependency) error
}
