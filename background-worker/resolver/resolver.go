package resolver

import (
	"github.com/oss-health/background-worker/internal/models"
)

type Resolver interface {
	CanResolve(dependency models.Dependency) bool
	Resolve(dependency models.Dependency) (string, error)
}
