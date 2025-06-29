package repository

import (
	"github.com/stretchr/testify/mock"

	"github.com/oss-health/background-worker/internal/dependency"
)

type MockExtractor struct {
	mock.Mock
}

func (m *MockExtractor) ExtractDependencies(path string) ([]dependency.DependencyVersionPair, error) {
	args := m.Called(path)
	var deps []dependency.DependencyVersionPair

	if args.Get(0) != nil {
		deps = args.Get(0).([]dependency.DependencyVersionPair)
	}

	return deps, args.Error(1)
}
