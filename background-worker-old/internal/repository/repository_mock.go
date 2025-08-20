package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetRepositoriesForDay(ctx context.Context, day int, hour int) ([]Repository, error) {
	args := m.Called(ctx, day, hour)
	return args.Get(0).([]Repository), args.Error(1)
}

func (m *MockRepository) MarkFailed(ctx context.Context, id int, reason string) error {
	args := m.Called(ctx, id, reason)
	return args.Error(0)
}

func (m *MockRepository) MarkScanned(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
