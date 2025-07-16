package resolver

import (
	"context"

	"github.com/oss-health/background-worker/internal/db"
)

type Resolver interface {
	CanResolve(dependency db.Dependency) bool
	Resolve(ctx context.Context, dependency db.Dependency) (string, error)
}
