package dependency

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) WaitUntilAllowed(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}
