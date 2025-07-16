package resolver

import (
	"context"
	"fmt"

	"github.com/oss-health/background-worker/internal/db"
)

type Chain interface {
	Resolve(ctx context.Context, dependency db.Dependency) (string, error)
}

type ResolverChain struct {
	resolvers []Resolver
}

func NewResolverChain(resolvers ...Resolver) *ResolverChain {
	return &ResolverChain{
		resolvers: resolvers,
	}
}

func (rc *ResolverChain) Resolve(ctx context.Context, dependency db.Dependency) (string, error) {
	for _, resolver := range rc.resolvers {
		if resolver.CanResolve(dependency) {
			return resolver.Resolve(ctx, dependency)
		}
	}
	return "", fmt.Errorf("no resolver found for ecosystem: %s", dependency.Ecosystem)
}
