package resolver

import (
	"context"
	"fmt"

	"github.com/oss-health/background-worker/internal/db"
)

type PyPIResolver struct{}

func (r *PyPIResolver) CanResolve(dependency db.Dependency) bool {
	return dependency.Ecosystem == "pypi"
}

func (r *PyPIResolver) Resolve(ctx context.Context, dependency db.Dependency) (string, error) {
	return fmt.Sprintf("https://pypi.org/project/%s/", dependency.Name), nil

}
