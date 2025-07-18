package resolver_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies/resolver"
	"golang.org/x/time/rate"
)

// MockResolver implements the Resolver interface for testing.
type MockResolver struct {
	resolveCalled bool
	canResolve    bool
	returnURL     string
	returnErr     error
}

func (m *MockResolver) CanResolve(dep db.Dependency) bool {
	return m.canResolve
}

func (m *MockResolver) Resolve(ctx context.Context, dep db.Dependency) (string, error) {
	m.resolveCalled = true
	return m.returnURL, m.returnErr
}

func TestRateLimitedResolver_CanResolve(t *testing.T) {
	mock := &MockResolver{canResolve: true}
	r := resolver.NewRateLimitedResolver(mock, 1, 1, "test")

	result := r.CanResolve(db.Dependency{Ecosystem: "npm"})
	if !result {
		t.Errorf("expected true from CanResolve, got false")
	}
}

func TestRateLimitedResolver_Resolve_Success(t *testing.T) {
	mock := &MockResolver{
		canResolve: true,
		returnURL:  "https://example.com",
		returnErr:  nil,
	}

	r := resolver.NewRateLimitedResolver(mock, 10, 1, "test")
	ctx := context.Background()
	dep := db.Dependency{Name: "lodash", Ecosystem: "npm"}

	url, err := r.Resolve(ctx, dep)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", url)
	}
	if !mock.resolveCalled {
		t.Errorf("expected Resolve to be called on inner resolver")
	}
}

func TestRateLimitedResolver_RespectsRateLimit(t *testing.T) {
	mock := &MockResolver{
		canResolve: true,
		returnURL:  "https://example.com",
		returnErr:  nil,
	}

	// Limit: 1 req/sec, burst = 1
	r := resolver.NewRateLimitedResolver(mock, 1, 1, "test")
	ctx := context.Background()
	dep := db.Dependency{Name: "lodash", Ecosystem: "npm"}

	// First request should succeed immediately
	start := time.Now()
	_, err := r.Resolve(ctx, dep)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	// Second request should be delayed by ~1s
	_, err = r.Resolve(ctx, dep)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}
	if elapsed < 1*time.Second {
		t.Errorf("rate limiter did not enforce delay: only waited %v", elapsed)
	}
}

func TestRateLimitedResolver_HandlesContextCancel(t *testing.T) {
	mock := &MockResolver{
		canResolve: true,
		returnURL:  "",
		returnErr:  nil,
	}

	r := resolver.NewRateLimitedResolver(mock, rate.Limit(0.1), 1, "test") // 1 every 10s

	// Use up the initial burst token so the next call must wait
	ctxImmediate := context.Background()
	_, _ = r.Resolve(ctxImmediate, db.Dependency{Name: "initial", Ecosystem: "npm"})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := r.Resolve(ctx, db.Dependency{Name: "lodash", Ecosystem: "npm"})
	if err == nil {
		t.Errorf("expected error due to context timeout, got nil")
	} else if !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "context deadline") {
		t.Errorf("expected rate limit context deadline error, got %v", err)
	}
}
