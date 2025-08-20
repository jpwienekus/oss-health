package repository

import (
	"github.com/stretchr/testify/mock"
)

type MockCloner struct {
	mock.Mock
}

func (m *MockCloner) CloneRepository(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

