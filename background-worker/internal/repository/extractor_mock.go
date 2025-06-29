package repository

import (
	"github.com/stretchr/testify/mock"

	"github.com/oss-health/background-worker/internal/repository/parsers"
)

type MockExtractor struct {
	mock.Mock
}

func (m *MockExtractor) ExtractDependencies(path string) ([]parsers.DependencyParsed, error) {
	args := m.Called(path)
	var deps []parsers.DependencyParsed

	if args.Get(0) != nil {
		deps = args.Get(0).([]parsers.DependencyParsed)
	}

	return deps, args.Error(1)
}
