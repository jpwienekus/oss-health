package fetcher_test

import (
	"context"
	"testing"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/pkg/fetcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResolvePendingDependencies_Success(t *testing.T) {
	ctx := context.Background()

	mockDB := new(MockDB)
	mockLimiter := new(MockRateLimiter)

	db.GetPendingDependencies = mockDB.GetPendingDependencies
	db.UpsertGithubURLs = mockDB.UpsertGithubURLs
	db.BatchUpdateDependencies = mockDB.BatchUpdateDependencies
	db.MarkDependenciesAsFailed = mockDB.MarkDependenciesAsFailed

	mockLimiter.On("WaitUntilAllowed", mock.Anything, "npm").Return(nil)
	mockResolvers := map[string]func(ctx context.Context, name string) (string, error){
		"npm": MockResolver,
	}

	dependencies := []db.Dependency{
		{ID: 1, Name: "package-one", Ecosystem: "npm"},
	}

	mockDB.On("GetPendingDependencies", ctx, 10, 0, "npm").Return(dependencies, nil)
	mockDB.On("UpsertGithubURLs", ctx, mock.Anything).Return(map[string]int64{
		"https://github.com/package-one": 123,
	}, nil)
	mockDB.On("BatchUpdateDependencies", ctx, dependencies, mock.Anything, mock.Anything).Return(nil)

	err := fetcher.ResolvePendingDependencies(ctx, 10, 0, "npm", mockLimiter, mockResolvers)
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockLimiter.AssertExpectations(t)
}

func TestResolvePendingDependencies_ResolverFailure(t *testing.T) {
	ctx := context.Background()

	mockDB := new(MockDB)
	mockLimiter := new(MockRateLimiter)

	db.GetPendingDependencies = mockDB.GetPendingDependencies
	db.MarkDependenciesAsFailed = mockDB.MarkDependenciesAsFailed

	mockLimiter.On("WaitUntilAllowed", mock.Anything, "npm").Return(nil)
	mockResolvers := map[string]func(ctx context.Context, name string) (string, error){
		"npm": MockResolver,
	}

	dependencies := []db.Dependency{
		{ID: 1, Name: "fail", Ecosystem: "npm"},
	}

	mockDB.On("GetPendingDependencies", ctx, 10, 0, "npm").Return(dependencies, nil)
	mockDB.On("MarkDependenciesAsFailed", ctx, mock.MatchedBy(func(failures map[int64]string) bool {
		_, exists := failures[1]
		return exists
	})).Return(nil)

	err := fetcher.ResolvePendingDependencies(ctx, 10, 0, "npm", mockLimiter, mockResolvers)
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockLimiter.AssertExpectations(t)
}


func TestResolvePendingDependencies_UnsupportedEcosystem(t *testing.T) {
	ctx := context.Background()

	mockDB := new(MockDB)
	mockLimiter := new(MockRateLimiter)

	db.GetPendingDependencies = mockDB.GetPendingDependencies
	db.MarkDependenciesAsFailed = mockDB.MarkDependenciesAsFailed

	mockResolvers := map[string]func(ctx context.Context, name string) (string, error){
		"npm": MockResolver,
	}
	dependencies := []db.Dependency{
		{ID: 1, Name: "foo", Ecosystem: "unknown"},
	}

	mockDB.On("GetPendingDependencies", ctx, 10, 0, "unknown").Return(dependencies, nil)
	mockDB.On("MarkDependenciesAsFailed", ctx, mock.MatchedBy(func(failures map[int64]string) bool {
		_, exists := failures[1]
		return exists
	})).Return(nil)

	err := fetcher.ResolvePendingDependencies(ctx, 10, 0, "unknown", mockLimiter, mockResolvers)
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockLimiter.AssertExpectations(t)
}
