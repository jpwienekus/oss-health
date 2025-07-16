package dependencies_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

// Mocks
type MockDependencyRepository struct {
	mock.Mock
}

func (m *MockDependencyRepository) GetExistingDependencies(ctx context.Context, deps []parsers.Dependency) ([]db.Dependency, error) {
	args := m.Called(ctx, deps)
	return args.Get(0).([]db.Dependency), args.Error(1)
}

func (m *MockDependencyRepository) InsertDependencies(ctx context.Context, deps []parsers.Dependency) error {
	args := m.Called(ctx, deps)
	return args.Error(0)
}

func (m *MockDependencyRepository) GetExistingVersionsForDependencies(ctx context.Context, deps []parsers.Dependency) ([]db.DependencyVersion, error) {
	args := m.Called(ctx, deps)
	return args.Get(0).([]db.DependencyVersion), args.Error(1)
}

func (m *MockDependencyRepository) InsertDependencyVersions(ctx context.Context, versions []db.DependencyVersion) error {
	args := m.Called(ctx, versions)
	return args.Error(0)
}

func (m *MockDependencyRepository) GetDependenciesPendingUrlResolution(ctx context.Context, repositoryId int) ([]db.Dependency, error) {
	args := m.Called(ctx, repositoryId)
	return args.Get(0).([]db.Dependency), args.Error(1)
}

func (m *MockDependencyRepository) InsertDependencyRepository(ctx context.Context, dependencyUrlMap map[int]string) error {
	args := m.Called(ctx, dependencyUrlMap)
	return args.Error(0)
}

func (m *MockDependencyRepository) GetDependencyRepositoryIds(ctx context.Context, dependencyUrlMap map[int]string) (map[string]int, error) {
	args := m.Called(ctx, dependencyUrlMap)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockDependencyRepository) UpdateDependencyRepositoryIDs(ctx context.Context, depRepoMap map[int]int) error {
	args := m.Called(ctx, depRepoMap)
	return args.Error(0)
}

func (m *MockDependencyRepository) MarkDependenciesAsFailed(ctx context.Context, failureReasons map[int]string) error {
	args := m.Called(ctx, failureReasons)
	return args.Error(0)
}

type MockResolverChain struct {
	mock.Mock
}

func (m *MockResolverChain) Resolve(ctx context.Context, dependency db.Dependency) (string, error) {
	args := m.Called(ctx, dependency)
	return args.String(0), args.Error(1)
}

// Tests
func TestInsertDependenciesThatDoNotExist(t *testing.T) {
	mockRepo := new(MockDependencyRepository)
	mockChain := new(MockResolverChain)
	svc := dependencies.NewDependencyService(mockRepo, mockChain)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm"},
		{Name: "express", Ecosystem: "npm"},
	}
	existing := []db.Dependency{
		{Name: "axios", Ecosystem: "npm", ID: 1},
	}
	expectedToInsert := []parsers.Dependency{
		{Name: "express", Ecosystem: "npm"},
	}

	mockRepo.On("GetExistingDependencies", ctx, input).Return(existing, nil)
	mockRepo.On("InsertDependencies", ctx, expectedToInsert).Return(nil)

	err := svc.InsertDependenciesThatDoNotExist(ctx, input, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInsertVersionsThatDoNotExist(t *testing.T) {
	mockRepo := new(MockDependencyRepository)
	mockChain := new(MockResolverChain)
	svc := dependencies.NewDependencyService(mockRepo, mockChain)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm", Version: "1.2.3"},
	}
	existingDeps := []db.Dependency{
		{ID: 1, Name: "axios", Ecosystem: "npm"},
	}
	existingVersions := []db.DependencyVersion{}

	expectedInsert := []db.DependencyVersion{
		{Version: "1.2.3", DependencyId: 1},
	}

	mockRepo.On("GetExistingDependencies", ctx, input).Return(existingDeps, nil)
	mockRepo.On("GetExistingVersionsForDependencies", ctx, input).Return(existingVersions, nil)
	mockRepo.On("InsertDependencyVersions", ctx, expectedInsert).Return(nil)

	err := svc.InsertVersionsThatDoNotExist(ctx, input, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetExistingVersionsForDependencies_Success(t *testing.T) {
	mockRepo := new(MockDependencyRepository)
	mockChain := new(MockResolverChain)
	svc := dependencies.NewDependencyService(mockRepo, mockChain)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm", Version: "1.0.0"},
	}
	expected := []db.DependencyVersion{
		{ID: 1, Version: "1.0.0", DependencyId: 1, DependencyName: "axios", DependencyEcosystem: "npm"},
	}

	mockRepo.On("GetExistingVersionsForDependencies", ctx, input).Return(expected, nil)

	result, err := svc.GetExistingVersionsForDependencies(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetExistingVersionsForDependencies_Empty(t *testing.T) {
	mockRepo := new(MockDependencyRepository)
	mockChain := new(MockResolverChain)
	svc := dependencies.NewDependencyService(mockRepo, mockChain)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm", Version: "1.0.0"},
	}

	mockRepo.On("GetExistingVersionsForDependencies", ctx, input).Return([]db.DependencyVersion{}, nil)

	_, err := svc.GetExistingVersionsForDependencies(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no existing versions found")
	mockRepo.AssertExpectations(t)
}
