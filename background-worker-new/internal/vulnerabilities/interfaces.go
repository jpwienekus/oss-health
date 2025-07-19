package vulnerabilities

import (
	"context"

	"github.com/oss-health/background-worker/internal/db"
)

type VulnerabilityRepositoryInterface interface {
	InsertVulnerabilities(ctx context.Context, dependencyVersionVulnerabilities map[int][]Vulnerability) ([]db.Vulnerability, error)
	InsertVersionVulnerability(ctx context.Context, versionVulnerabilities map[int]int) error
}
