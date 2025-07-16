package resolver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oss-health/background-worker/internal/db"
	"golang.org/x/time/rate"
)

type RateLimitedResolver struct {
	resolver Resolver
	limiter  *rate.Limiter

	mu          sync.Mutex
	requests    int
	windowStart time.Time
}

func NewRateLimitedResolver(resolver Resolver, r rate.Limit, b int) *RateLimitedResolver {
	return &RateLimitedResolver{
		resolver: resolver,
		limiter:  rate.NewLimiter(r, b),
	}
}

func (r *RateLimitedResolver) CanResolve(dependency db.Dependency) bool {
	return r.resolver.CanResolve(dependency)
}

func (r *RateLimitedResolver) Resolve(ctx context.Context, dependency db.Dependency) (string, error) {
	if !r.limiter.Allow() {
		fmt.Printf("Rate limit reached for dependency: %+v. Waiting...\n", dependency)
	}

	if err := r.limiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limit exceeded: %w", err)
	}

	r.trackRequest()
	return r.resolver.Resolve(ctx, dependency)
}

func (r *RateLimitedResolver) trackRequest() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if r.windowStart.IsZero() || now.Sub(r.windowStart) >= time.Second {
		r.windowStart = now
		r.requests = 0
	}
	r.requests++
}

func (r *RateLimitedResolver) requestsThisSecond() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.requests
}

func (r *RateLimitedResolver) StartLogging(name string) {
	go func() {
		for range time.Tick(time.Second) {
			fmt.Printf("[%s] Requests this second: %d\n", name, r.requestsThisSecond())
		}
	}()
}
