package resolver

import (
	"context"
	"fmt"

	"github.com/oss-health/background-worker/internal/models"
	"golang.org/x/time/rate"
)

type RateLimitedResolver struct {
	resolver Resolver
	limiter *rate.Limiter
}

func NewRateLimitedResolver(resolver Resolver, r rate.Limit, b int) *RateLimitedResolver {
	return &RateLimitedResolver{
		resolver: resolver,
		limiter:  rate.NewLimiter(r, b),
	}
}


func (r *RateLimitedResolver) CanResolve(dependency models.Dependency) bool {
	return r.resolver.CanResolve(dependency)
}

func (r *RateLimitedResolver) Resolve(dependency models.Dependency) (string, error) {
	if err := r.limiter.Wait(context.Background()); err != nil {
		return "", fmt.Errorf("rate limit exceeded: %w", err)
	}

	return r.resolver.Resolve(dependency)
}
