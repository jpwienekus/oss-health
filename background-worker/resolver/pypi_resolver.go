package resolver

import (
	"fmt"

	"github.com/oss-health/background-worker/internal/models"
)

type PyPIResolver struct{}

func (r *PyPIResolver) CanResolve(dependency models.Dependency) bool {
	return dependency.Ecosystem == "pypi"
}

func (r *PyPIResolver) Resolve(dependency models.Dependency) (string, error) {
	return fmt.Sprintf("https://pypi.org/project/%s/", dependency.Name), nil

}
