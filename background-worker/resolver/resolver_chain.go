package resolver

import (
	"fmt"

	"github.com/oss-health/background-worker/internal/models"
)

type Chain interface {
	Resolve(models.Dependency) (string, error)
}

type ResolverChain struct {
	resolvers []Resolver
}

func NewResolverChain(resolvers ...Resolver) *ResolverChain {
	return &ResolverChain{
		resolvers: resolvers,
	}
}

func (rc *ResolverChain) Resolve(dependency models.Dependency) (string, error) {
	for _, resolver := range rc.resolvers {
		if resolver.CanResolve(dependency) {
			return resolver.Resolve(dependency)
		}
	}
	return "", fmt.Errorf("no resolver found for ecosystem: %s", dependency.Ecosystem)
}
