package utils

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type DefaultRateLimiter struct {
	limiters map[string]*rate.Limiter
	mutex    sync.RWMutex
}

func NewDefaultRateLimiter() *DefaultRateLimiter {
	return &DefaultRateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (r *DefaultRateLimiter) RegisterLimiter(ecosystem string, limit rate.Limit, burst int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.limiters[ecosystem] = rate.NewLimiter(limit, burst)
}

func (r *DefaultRateLimiter) WaitUntilAllowed(ctx context.Context, ecosystem string) error {
	r.mutex.RLock()
	limiter, ok := r.limiters[ecosystem]
	r.mutex.RUnlock()

	if !ok {
		return nil
	}
	return limiter.Wait(ctx)
}
// package utils
//
// import (
// 	"context"
// 	"sync"
//
// 	"golang.org/x/time/rate"
// )
//
// var (
// 	limiters = make(map[string]*rate.Limiter)
// 	mutex    sync.Mutex
// )
//
// func RegisterLimiter(ecosystem string, limit rate.Limit, burst int) {
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	limiters[ecosystem] = rate.NewLimiter(limit, burst)
// }
//
// var WaitUntilAllowed = func(ctx context.Context, ecosystem string) error {
// 	mutex.Lock()
// 	limiter, ok := limiters[ecosystem]
// 	mutex.Unlock()
//
// 	if !ok {
// 		return nil
// 	}
//
// 	return limiter.Wait(ctx)
// }
