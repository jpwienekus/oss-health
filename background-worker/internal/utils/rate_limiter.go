package utils

import "context"

type RateLimiter interface {
	WaitUntilAllowed(ctx context.Context, ecosystem string) error
}
