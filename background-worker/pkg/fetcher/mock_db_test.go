package fetcher_test

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/oss-health/background-worker/internal/db"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetPendingDependencies(ctx context.Context, batchSize, offset int, ecosystem string) ([]db.Dependency, error) {
	args := m.Called(ctx, batchSize, offset, ecosystem)
	return args.Get(0).([]db.Dependency), args.Error(1)
}
func (m *MockDB) UpsertGithubURLs(ctx context.Context, urls []string) (map[string]int64, error) {
	args := m.Called(ctx, urls)
	return args.Get(0).(map[string]int64), args.Error(1)
}
func (m *MockDB) BatchUpdateDependencies(ctx context.Context, deps []db.Dependency, urlToID map[string]int64, resolved map[int64]string) error {
	args := m.Called(ctx, deps, urlToID, resolved)
	return args.Error(0)
}
func (m *MockDB) MarkDependenciesAsFailed(ctx context.Context, failures map[int64]string) error {
	args := m.Called(ctx, failures)
	return args.Error(0)
}

// Mock resolver
func MockResolver(ctx context.Context, name string) (string, error) {
	if name == "fail" {
		return "", errors.New("resolver failed")
	}
	return "https://github.com/" + name, nil
}

// Mock rate limiter (disable real rate limiting for test)
type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) WaitUntilAllowed(ctx context.Context, ecosystem string) error {
	args := m.Called(ctx, ecosystem)
	return args.Error(0)
}
