package resolver

import (
	"context"

	"github.com/oss-health/background-worker/internal/models"
)

type Resolver interface {
	CanResolve(dependency models.Dependency) bool
	Resolve(ctx context.Context, dependency models.Dependency) (string, error)
}
