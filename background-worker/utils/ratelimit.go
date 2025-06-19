package utils

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// Register a rate limiter for an ecosystem (e.g. "npm", "pypi")
func RegisterLimiter(ecosystem string, limit rate.Limit, burst int) {
	mu.Lock()
	defer mu.Unlock()
	limiters[ecosystem] = rate.NewLimiter(limit, burst)
}

// Wait until a request is allowed for this ecosystem
func WaitUntilAllowed(ctx context.Context, ecosystem string) error {
	mu.Lock()
	limiter, ok := limiters[ecosystem]
	mu.Unlock()

	if !ok {
		return nil // No limiter set = no limit
	}

	return limiter.Wait(ctx)
}
