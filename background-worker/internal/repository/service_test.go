package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/repository"
)

func TestProcessRepository_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockRepository)
	mockDependencyRepo := new(dependency.MockDependencyRepository)
	mockCloner := new(repository.MockCloner)
	mockExtractor := new(repository.MockExtractor)

	service := repository.NewRepositoryService(mockRepo, mockDependencyRepo, mockCloner, mockExtractor)

	repo := repository.Repository{
		ID:  1,
		URL: "https://github.com/example/repo",
	}

	mockCloner.On("CloneRepository", repo.URL).Return("/tmp/fakepath", nil)
	mockRepo.On("MarkScanned", ctx, repo.ID).Return(nil)
	mockExtractor.On("ExtractDependencies", "/tmp/fakepath").Return([]dependency.DependencyVersionPair{
		{Name: "dep1", Version: "1.0"},
	}, nil)

	deps, err := service.ProcessRepository(ctx, repo)

	assert.NoError(t, err)
	assert.Len(t, deps, 1)
	assert.Equal(t, "dep1", deps[0].Name)

	mockCloner.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockExtractor.AssertExpectations(t)
}

func TestProcessRepository_CloneFails(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockRepository)
	mockDependencyRepo := new(dependency.MockDependencyRepository)
	mockCloner := new(repository.MockCloner)
	mockExtractor := new(repository.MockExtractor)

	service := repository.NewRepositoryService(mockRepo, mockDependencyRepo, mockCloner, mockExtractor)

	repo := repository.Repository{ID: 1, URL: "invalid-url"}

	mockCloner.On("CloneRepository", repo.URL).Return("", errors.New("clone failed"))
	mockRepo.On("MarkFailed", ctx, repo.ID, "clone failed").Return(nil)

	deps, err := service.ProcessRepository(ctx, repo)

	assert.Error(t, err)
	assert.Nil(t, deps)

	mockCloner.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestProcessRepository_ExtractFails(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockRepository)
	mockDependencyRepo := new(dependency.MockDependencyRepository)
	mockCloner := new(repository.MockCloner)
	mockExtractor := new(repository.MockExtractor)

	service := repository.NewRepositoryService(mockRepo, mockDependencyRepo, mockCloner, mockExtractor)

	repo := repository.Repository{ID: 1, URL: "https://github.com/example/repo"}

	mockCloner.On("CloneRepository", repo.URL).Return("/tmp/fake", nil)
	mockRepo.On("MarkScanned", ctx, repo.ID).Return(nil)
	mockExtractor.On("ExtractDependencies", "/tmp/fake").Return(nil, errors.New("parser error"))

	deps, err := service.ProcessRepository(ctx, repo)

	assert.Error(t, err)
	assert.Nil(t, deps)

	mockCloner.AssertExpectations(t)
	mockExtractor.AssertExpectations(t)
}

func TestCloneAndParse_LogsError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockRepository)
	mockDependencyRepo := new(dependency.MockDependencyRepository)
	mockCloner := new(repository.MockCloner)
	mockExtractor := new(repository.MockExtractor)

	service := repository.NewRepositoryService(mockRepo, mockDependencyRepo, mockCloner, mockExtractor)

	repo := repository.Repository{ID: 1, URL: "bad-url"}

	mockCloner.On("CloneRepository", repo.URL).Return("", errors.New("mock failure"))
	mockRepo.On("MarkFailed", ctx, repo.ID, "mock failure").Return(nil)

	service.CloneAndParse(ctx, repo)

	mockCloner.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestRunDailyScan_CallsGetRepositories(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(repository.MockRepository)
	mockDependencyRepo := new(dependency.MockDependencyRepository)
	mockCloner := new(repository.MockCloner)
	mockExtractor := new(repository.MockExtractor)

	service := repository.NewRepositoryService(mockRepo, mockDependencyRepo, mockCloner, mockExtractor)

	repos := []repository.Repository{
		{ID: 1, URL: "https://github.com/a"},
		{ID: 2, URL: "https://github.com/b"},
	}

	mockRepo.On("GetRepositoriesForDay", ctx, 5, 15).Return(repos, nil)
	mockCloner.On("CloneRepository", mock.Anything).Return("/tmp/fake", nil)
	mockRepo.On("MarkScanned", ctx, 1).Return(nil)
	mockRepo.On("MarkScanned", ctx, 2).Return(nil)
	mockExtractor.On("ExtractDependencies", mock.Anything).Return([]dependency.DependencyVersionPair{}, nil)
  mockDependencyRepo.On("ReplaceRepositoryDependencyVersions", mock.Anything, mock.Anything, mock.Anything).Return([]dependency.DependencyVersionResult{}, nil)

	service.RunDailyScan(ctx, 5, 15)

	mockRepo.AssertExpectations(t)
}
