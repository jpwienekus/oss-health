package resolver

import (
	"fmt"

	"github.com/oss-health/background-worker/internal/models"
)

type NPMResolver struct{}

func (r *NPMResolver) CanResolve(dependency models.Dependency) bool {
	return dependency.Ecosystem == "npm"
}

func (r *NPMResolver) Resolve(dependency models.Dependency) (string, error) {
	return fmt.Sprintf("https://www.npmjs.com/package/%s", dependency.Name), nil
}
