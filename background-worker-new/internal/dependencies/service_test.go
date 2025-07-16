package dependencies_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/oss-health/background-worker/internal/dependencies"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

// Mocks
type MockDependencyRepository struct {
	mock.Mock
}

func (m *MockDependencyRepository) GetExistingDependencies(ctx context.Context, deps []parsers.Dependency) ([]dependencies.Dependency, error) {
	args := m.Called(ctx, deps)
	return args.Get(0).([]dependencies.Dependency), args.Error(1)
}

func (m *MockDependencyRepository) InsertDependencies(ctx context.Context, deps []parsers.Dependency) error {
	args := m.Called(ctx, deps)
	return args.Error(0)
}

func (m *MockDependencyRepository) GetExistingVersionsForDependencies(ctx context.Context, deps []parsers.Dependency) ([]dependencies.DependencyVersion, error) {
	args := m.Called(ctx, deps)
	return args.Get(0).([]dependencies.DependencyVersion), args.Error(1)
}

func (m *MockDependencyRepository) InsertDependencyVersions(ctx context.Context, versions []dependencies.DependencyVersion) error {
	args := m.Called(ctx, versions)
	return args.Error(0)
}

// Tests
func TestInsertDependenciesThatDoNotExist(t *testing.T) {
	mockRepo := new(MockDependencyRepository)
	svc := dependencies.NewDependencyService(mockRepo)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm"},
		{Name: "express", Ecosystem: "npm"},
	}
	existing := []dependencies.Dependency{
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
	svc := dependencies.NewDependencyService(mockRepo)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm", Version: "1.2.3"},
	}
	existingDeps := []dependencies.Dependency{
		{ID: 1, Name: "axios", Ecosystem: "npm"},
	}
	existingVersions := []dependencies.DependencyVersion{}

	expectedInsert := []dependencies.DependencyVersion{
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
	svc := dependencies.NewDependencyService(mockRepo)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm", Version: "1.0.0"},
	}
	expected := []dependencies.DependencyVersion{
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
	svc := dependencies.NewDependencyService(mockRepo)

	ctx := context.Background()
	input := []parsers.Dependency{
		{Name: "axios", Ecosystem: "npm", Version: "1.0.0"},
	}

	mockRepo.On("GetExistingVersionsForDependencies", ctx, input).Return([]dependencies.DependencyVersion{}, nil)

	_, err := svc.GetExistingVersionsForDependencies(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no existing versions found")
	mockRepo.AssertExpectations(t)
}
