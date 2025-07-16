package repo_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repo"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

// Mocks
type MockRepoRepository struct {
	mock.Mock
}

func (m *MockRepoRepository) GetReposForSchedule(ctx context.Context, day int, hour int) ([]db.Repository, error) {
	args := m.Called(ctx, day, hour)
	return args.Get(0).([]db.Repository), args.Error(1)
}

func (m *MockRepoRepository) MarkScanned(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockRepoRepository) MarkFailed(ctx context.Context, id int, reason string) error {
	return m.Called(ctx, id, reason).Error(0)
}

func (m *MockRepoRepository) DeleteRepositoryDependencyVersion(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockRepoRepository) InsertRepositoryDependencyVersion(ctx context.Context, id int, versions []db.DependencyVersion) error {
	return m.Called(ctx, id, versions).Error(0)
}

type MockDependencyService struct {
	mock.Mock
}

func (m *MockDependencyService) InsertDependenciesThatDoNotExist(ctx context.Context, deps []parsers.Dependency, workerId int) error {
	return m.Called(ctx, deps, workerId).Error(0)
}

func (m *MockDependencyService) InsertVersionsThatDoNotExist(ctx context.Context, deps []parsers.Dependency, workerId int) error {
	return m.Called(ctx, deps, workerId).Error(0)
}

func (m *MockDependencyService) GetExistingVersionsForDependencies(ctx context.Context, deps []parsers.Dependency) ([]db.DependencyVersion, error) {
	args := m.Called(ctx, deps)
	return args.Get(0).([]db.DependencyVersion), args.Error(1)
}

// Tests
func TestCloneRepository_Success(t *testing.T) {
	mockRepo := new(MockRepoRepository)
	mockDep := new(MockDependencyService)

	svc := repo.NewRepositoryService(mockRepo, mockDep, nil)

	// override execCommand to simulate git clone success
	repo.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "cloned")
	}

	defer func() { repo.ExecCommand = exec.Command }() // restore original

	err := svc.CloneRepository("http://fake.url/repo.git", "/fake/path")
	assert.NoError(t, err)
}

func TestReplaceRepositoryDependencyVersions_Success(t *testing.T) {
	mockRepo := new(MockRepoRepository)
	mockDep := new(MockDependencyService)

	repoID := 1
	deps := []parsers.Dependency{
		{Name: "lodash", Ecosystem: "npm", Version: "1.0.0"},
	}

	mockDep.On("InsertDependenciesThatDoNotExist", mock.Anything, deps, 1).Return(nil)
	mockDep.On("InsertVersionsThatDoNotExist", mock.Anything, deps, 1).Return(nil)

	mockRepo.On("DeleteRepositoryDependencyVersion", mock.Anything, repoID).Return(nil)
	mockDep.On("GetExistingVersionsForDependencies", mock.Anything, deps).
		Return([]db.DependencyVersion{{ID: 1, Version: "1.0.0", DependencyId: 2}}, nil)
	mockRepo.On("InsertRepositoryDependencyVersion", mock.Anything, repoID, mock.Anything).Return(nil)

	svc := repo.NewRepositoryService(mockRepo, mockDep, nil)

	err := svc.ReplaceRepositoryDependencyVersions(context.Background(), repoID, deps, 1)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockDep.AssertExpectations(t)
}

func TestRunDailyScan_InvalidParallelEnv(t *testing.T) {
	mockRepo := new(MockRepoRepository)

	os.Setenv("MAX_PARALLEL", "not-a-number")
	defer os.Unsetenv("MAX_PARALLEL")

	svc := repo.NewRepositoryService(mockRepo, nil, nil)
	mockRepo.On("GetReposForSchedule", mock.Anything, 1, 10).Return([]db.Repository{}, nil)

	err := svc.RunDailyScan(context.Background(), 1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid MAX_PARALLEL")
}
