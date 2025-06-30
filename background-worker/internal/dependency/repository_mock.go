package dependency

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDependencyRepository struct {
    mock.Mock
}
func (m *MockDependencyRepository) GetDependenciesPendingUrlResolution(ctx context.Context, batchSize, offset int, eco string) ([]Dependency, error) {
    args := m.Called(ctx, batchSize, offset, eco)
    return args.Get(0).([]Dependency), args.Error(1)
}
func (m *MockDependencyRepository) UpsertGithubURLs(ctx context.Context, urls []string) (map[string]int64, error) {
    args := m.Called(ctx, urls)
    return args.Get(0).(map[string]int64), args.Error(1)
}
func (m *MockDependencyRepository) BatchUpdateDependencies(ctx context.Context, deps []Dependency, urlToID map[string]int64, resolvedURLs map[int64]string) error {
    args := m.Called(ctx, deps, urlToID, resolvedURLs)
    return args.Error(0)
}
func (m *MockDependencyRepository) MarkDependenciesAsFailed(ctx context.Context, failures map[int64]string) error {
    args := m.Called(ctx, failures)
    return args.Error(0)
}

func (m *MockDependencyRepository) ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryID int, pairs []DependencyVersionPair) ([]DependencyVersionResult, error) {
    args := m.Called(ctx, repositoryID, pairs)
    return args.Get(0).([]DependencyVersionResult), args.Error(1)
}

