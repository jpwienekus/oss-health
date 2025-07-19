package vulnerabilities

import (
	"context"
	"errors"
	"testing"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repo/parsers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockDependencyRepo struct {
	mock.Mock
}

func (m *MockDependencyRepo) GetExistingDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.Dependency, error) {
	args := m.Called(ctx, dependencies)
	return args.Get(0).([]db.Dependency), args.Error(1)
}
func (m *MockDependencyRepo) GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.DependencyVersion, error) {
	args := m.Called(ctx, dependencies)
	return args.Get(0).([]db.DependencyVersion), args.Error(1)
}
func (m *MockDependencyRepo) InsertDependencies(ctx context.Context, dependencies []parsers.Dependency) error {
	args := m.Called(ctx, dependencies)
	return args.Error(0)
}
func (m *MockDependencyRepo) InsertDependencyVersions(ctx context.Context, versions []db.DependencyVersion) error {
	args := m.Called(ctx, versions)
	return args.Error(0)
}
func (m *MockDependencyRepo) GetDependenciesPendingUrlResolution(ctx context.Context, repoID int) ([]db.Dependency, error) {
	args := m.Called(ctx, repoID)
	return args.Get(0).([]db.Dependency), args.Error(1)
}
func (m *MockDependencyRepo) InsertDependencyRepository(ctx context.Context, dependencyUrlMap map[int]string) error {
	args := m.Called(ctx, dependencyUrlMap)
	return args.Error(0)
}
func (m *MockDependencyRepo) GetDependencyRepositoryIds(ctx context.Context, dependencyUrlMap map[int]string) (map[string]int, error) {
	args := m.Called(ctx, dependencyUrlMap)
	return args.Get(0).(map[string]int), args.Error(1)
}
func (m *MockDependencyRepo) UpdateDependencyRepositoryIDs(ctx context.Context, depRepoMap map[int]int) error {
	args := m.Called(ctx, depRepoMap)
	return args.Error(0)
}
func (m *MockDependencyRepo) MarkDependenciesAsFailed(ctx context.Context, failureReasons map[int]string) error {
	args := m.Called(ctx, failureReasons)
	return args.Error(0)
}
func (m *MockDependencyRepo) GetDependenciesPendingOSVScan(ctx context.Context) ([]db.DependencyVersion, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.DependencyVersion), args.Error(1)
}

type MockVulnerabilityRepo struct {
	mock.Mock
}

func (m *MockVulnerabilityRepo) InsertVulnerabilities(ctx context.Context, vulnerabilities map[int][]Vulnerability) ([]db.Vulnerability, error) {
	args := m.Called(ctx, vulnerabilities)
	return args.Get(0).([]db.Vulnerability), args.Error(1)
}

func (m *MockVulnerabilityRepo) InsertVersionVulnerability(ctx context.Context, versionVulns map[int]int) error {
	args := m.Called(ctx, versionVulns)
	return args.Error(0)
}

type MockOSVClient struct {
	mock.Mock
}

func (m *MockOSVClient) QueryBatch(ctx context.Context, deps []db.DependencyVersion) (map[int][]string, error) {
	args := m.Called(ctx, deps)
	return args.Get(0).(map[int][]string), args.Error(1)
}

func (m *MockOSVClient) GetVulnerabilityDetails(ctx context.Context, vulns map[int][]string) (map[int][]Vulnerability, error) {
	args := m.Called(ctx, vulns)
	return args.Get(0).(map[int][]Vulnerability), args.Error(1)
}

// --- Tests ---

func TestRunOSVScan_Success(t *testing.T) {
	ctx := context.Background()

	mockDepRepo := new(MockDependencyRepo)
	mockVulnRepo := new(MockVulnerabilityRepo)
	mockOSVClient := new(MockOSVClient)

	service := NewVulnerabilityService(mockVulnRepo, mockDepRepo, mockOSVClient)

	// Prepare mock data
	depVersions := []db.DependencyVersion{{ID: 1, Version: "1.0.0"}}
	osvVulns := map[int][]Vulnerability{
		1: {{ID: "OSV-1", Summary: "example vuln"}},
	}
	insertedVulns := []db.Vulnerability{{ID: 100, OSVID: "OSV-1"}}
	versionVulnMap := map[int]int{100: 1}

	// Expectations
	osvIDs := map[int][]string{
		1: {"OSV-1"},
	}
	// osvVulns := map[int][]Vulnerability{
	// 	1: {{ID: "OSV-1", Summary: "example vuln"}},
	// }
	mockDepRepo.On("GetDependenciesPendingOSVScan", ctx).Return([]db.DependencyVersion{
		{ID: 1, Version: "1.0.0"},
	}, nil)
	mockOSVClient.On("QueryBatch", ctx, depVersions).Return(osvIDs, nil)
	mockOSVClient.On("GetVulnerabilityDetails", ctx, osvIDs).Return(osvVulns, nil)
	mockVulnRepo.On("InsertVulnerabilities", ctx, osvVulns).Return(insertedVulns, nil)
	mockVulnRepo.On("InsertVersionVulnerability", ctx, versionVulnMap).Return(nil)

	// Run
	err := service.RunOSVScan(ctx)

	// Verify
	assert.NoError(t, err)
	mockDepRepo.AssertExpectations(t)
	mockVulnRepo.AssertExpectations(t)
	mockOSVClient.AssertExpectations(t)
}

func TestRunOSVScan_NoDeps(t *testing.T) {
	ctx := context.Background()
	mockDepRepo := new(MockDependencyRepo)
	mockVulnRepo := new(MockVulnerabilityRepo)
	mockOSVClient := new(MockOSVClient)

	service := NewVulnerabilityService(mockVulnRepo, mockDepRepo, mockOSVClient)

	mockDepRepo.On("GetDependenciesPendingOSVScan", ctx).Return([]db.DependencyVersion{}, nil)

	err := service.RunOSVScan(ctx)

	assert.ErrorContains(t, err, "no dependencies pending OSV scanning")
}

func TestRunOSVScan_ErrorInQueryBatch(t *testing.T) {
	ctx := context.Background()
	mockDepRepo := new(MockDependencyRepo)
	mockVulnRepo := new(MockVulnerabilityRepo)
	mockOSVClient := new(MockOSVClient)

	service := NewVulnerabilityService(mockVulnRepo, mockDepRepo, mockOSVClient)

	depVersions := []db.DependencyVersion{{ID: 1, Version: "1.0.0"}}
	mockDepRepo.On("GetDependenciesPendingOSVScan", ctx).Return(depVersions, nil)
	mockOSVClient.On("QueryBatch", ctx, depVersions).Return(make(map[int][]string), errors.New("query failed"))

	err := service.RunOSVScan(ctx)

	assert.ErrorContains(t, err, "failed to get OSV vulnerabilities")
}
