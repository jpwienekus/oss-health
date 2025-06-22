// package fetcher_test
//
// import (
// 	"context"
//
// 	"github.com/stretchr/testify/mock"
// 	"github.com/oss-health/background-worker/pkg/db"
// )
//
// type MockDB struct {
// 	mock.Mock
// }
//
// func (m *MockDB) GetPendingDependencies(ctx context.Context, batchSize, offset int, ecosystem string) ([]db.Dependency, error) {
// 	args := m.Called(ctx, batchSize, offset, ecosystem)
// 	return args.Get(0).([]db.Dependency), args.Error(1)
// }
