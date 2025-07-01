package dependency_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/oss-health/background-worker/internal/dependency"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResolvePendingDependencies_Success(t *testing.T) {
	ctx := context.Background()
	repository := new(dependency.MockDependencyRepository)
	limiter := new(dependency.MockRateLimiter)

	deps := []dependency.Dependency{
		{ID: 1, Name: "test", Ecosystem: "npm"},
	}

	repository.On("GetDependenciesPendingUrlResolution", ctx, 10, 0, "npm").Return(deps, nil)
	limiter.On("WaitUntilAllowed", ctx, "npm").Return(nil)
	repository.On("UpsertGithubURLs", ctx, map[int64]string{1: "https://github.com/test/repo"}).
		Return(map[int64]int64{100: 10}, nil)
	repository.On("BatchUpdateDependencies", ctx, map[int64]int64{100: 10}).Return(nil)

	service := dependency.NewDependencyService(repository, limiter, map[string]func(context.Context, string) (string, error){
		"npm": func(ctx context.Context, name string) (string, error) {
			return "https://github.com/test/repo", nil
		},
	})

	err := service.ResolvePendingDependencies(ctx, 10, 0, "npm")
	assert.NoError(t, err)
	repository.AssertExpectations(t)
	limiter.AssertExpectations(t)
}

func TestResolvePendingDependencies_EmptyList(t *testing.T) {
	ctx := context.Background()
	repository := new(dependency.MockDependencyRepository)
	limiter := new(dependency.MockRateLimiter)

	repository.On("GetDependenciesPendingUrlResolution", ctx, 10, 0, "npm").Return([]dependency.Dependency{}, nil)

	service := dependency.NewDependencyService(repository, limiter, nil)

	err := service.ResolvePendingDependencies(ctx, 10, 0, "npm")
	assert.NoError(t, err)
}

func TestResolvePendingDependencies_UnsupportedEcosystem(t *testing.T) {
	ctx := context.Background()
	repository := new(dependency.MockDependencyRepository)
	limiter := new(dependency.MockRateLimiter)

	deps := []dependency.Dependency{{ID: 1, Name: "abc", Ecosystem: "unknown"}}
	repository.On("GetDependenciesPendingUrlResolution", ctx, 10, 0, "unknown").Return(deps, nil)
	repository.On("MarkDependenciesAsFailed", ctx, map[int64]string{1: "unsupported ecosystem"}).Return(nil)

	service := dependency.NewDependencyService(repository, limiter, nil)

	err := service.ResolvePendingDependencies(ctx, 10, 0, "unknown")
	assert.NoError(t, err)
}

func TestResolvePendingDependencies_RateLimiterError(t *testing.T) {
	ctx := context.Background()
	repository := new(dependency.MockDependencyRepository)
	limiter := new(dependency.MockRateLimiter)

	deps := []dependency.Dependency{{ID: 1, Name: "rate", Ecosystem: "npm"}}
	repository.On("GetDependenciesPendingUrlResolution", ctx, 10, 0, "npm").Return(deps, nil)
	repository.On("MarkDependenciesAsFailed", ctx, mock.MatchedBy(func(failures map[int64]string) bool {
		return failures[1] == "rate limit: rate limited"
	})).Return(nil)
	limiter.On("WaitUntilAllowed", ctx, "npm").Return(errors.New("rate limited"))

	service := dependency.NewDependencyService(repository, limiter, map[string]func(context.Context, string) (string, error){
		"npm": func(ctx context.Context, name string) (string, error) {
			return "https://github.com/test/repo", nil
		},
	})

	err := service.ResolvePendingDependencies(ctx, 10, 0, "npm")
	assert.NoError(t, err)
}

func TestResolvePendingDependencies_ResolverError(t *testing.T) {
	ctx := context.Background()
	repository := new(dependency.MockDependencyRepository)
	limiter := new(dependency.MockRateLimiter)

	deps := []dependency.Dependency{{ID: 1, Name: "fail", Ecosystem: "npm"}}
	repository.On("GetDependenciesPendingUrlResolution", ctx, 10, 0, "npm").Return(deps, nil)
	limiter.On("WaitUntilAllowed", ctx, "npm").Return(nil)
	repository.On("MarkDependenciesAsFailed", ctx, mock.Anything).Return(nil)

	service := dependency.NewDependencyService(repository, limiter, map[string]func(context.Context, string) (string, error){
		"npm": func(ctx context.Context, name string) (string, error) {
			return "", errors.New("resolver failed")
		},
	})

	err := service.ResolvePendingDependencies(ctx, 10, 0, "npm")
	assert.NoError(t, err)
}

func TestResolvePendingDependencies_UpsertGithubURLsFails(t *testing.T) {
	ctx := context.Background()
	repository := new(dependency.MockDependencyRepository)
	limiter := new(dependency.MockRateLimiter)

	deps := []dependency.Dependency{{ID: 1, Name: "test", Ecosystem: "npm"}}
	repository.On("GetDependenciesPendingUrlResolution", ctx, 10, 0, "npm").Return(deps, nil)
	limiter.On("WaitUntilAllowed", ctx, "npm").Return(nil)
	repository.On("UpsertGithubURLs", ctx, mock.Anything).Return(map[int64]int64{}, errors.New("upsert failed"))

	service := dependency.NewDependencyService(repository, limiter, map[string]func(context.Context, string) (string, error){
		"npm": func(ctx context.Context, name string) (string, error) {
			return "https://github.com/test/repo", nil
		},
	})

	err := service.ResolvePendingDependencies(ctx, 10, 0, "npm")
	assert.ErrorContains(t, err, "upsert GitHub URLs: upsert failed")
}

func TestResolvePendingDependencies_Concurrency(t *testing.T) {
	ctx := context.Background()
	repository := new(dependency.MockDependencyRepository)
	limiter := new(dependency.MockRateLimiter)

	depCount := 100
	deps := make([]dependency.Dependency, depCount)

	for i := 0; i < depCount; i++ {
		deps[i] = dependency.Dependency{
			ID:        int64(i + 1),
			Name:      fmt.Sprintf("dep-%d", i),
			Ecosystem: "npm",
		}
	}

	repository.On("GetDependenciesPendingUrlResolution", ctx, depCount, 0, "npm").Return(deps, nil)

	for i := 0; i < depCount; i++ {
		limiter.On("WaitUntilAllowed", ctx, "npm").Return(nil)
	}

	upserted := make(map[int64]int64)

	for _, dep := range deps {
		upserted[dep.ID] = dep.ID + 1000
	}

	repository.On("UpsertGithubURLs", ctx, mock.Anything).Return(upserted, nil)

	repository.On("BatchUpdateDependencies", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	var mu sync.Mutex
	calls := 0

	resolver := func(ctx context.Context, name string) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		calls++
		return fmt.Sprintf("https://github.com/org/%s", name), nil
	}

	service := dependency.NewDependencyService(repository, limiter, map[string]func(context.Context, string) (string, error){
		"npm": resolver,
	})

	err := service.ResolvePendingDependencies(ctx, depCount, 0, "npm")
	assert.NoError(t, err)

	mu.Lock()
	assert.Equal(t, depCount, calls)
	mu.Unlock()

	repository.AssertExpectations(t)
	limiter.AssertExpectations(t)
}
