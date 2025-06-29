package repository

import (
	"github.com/stretchr/testify/mock"
)

type MockExtractor struct {
	mock.Mock
}

func (m *MockExtractor) ExtractDependencies(path string) ([]DependencyParsed, error) {
	args := m.Called(path)
	var deps []DependencyParsed

	if args.Get(0) != nil {
		deps = args.Get(0).([]DependencyParsed)
	}

	return deps, args.Error(1)
}
