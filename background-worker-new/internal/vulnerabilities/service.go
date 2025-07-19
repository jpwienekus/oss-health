package vulnerabilities

import (
	"context"
	"fmt"
	"log"

	"github.com/oss-health/background-worker/internal/dependencies"
)

type VulnerabilityService struct {
	repository           VulnerabilityRepositoryInterface
	dependencyRepository dependencies.DependencyRepositoryInterface
}

func NewVulnerabilityService(repository VulnerabilityRepositoryInterface, dependencyRepository dependencies.DependencyRepositoryInterface) *VulnerabilityService {
	return &VulnerabilityService{
		repository:           repository,
		dependencyRepository: dependencyRepository,
	}
}

// TODO: continue tests from here:
func (s *VulnerabilityService) RunOSVScan(ctx context.Context) error {
	dependencyVersions, err := s.dependencyRepository.GetDependenciesPendingOSVScan(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dependencies pending OSV scanning: %w", err)
	}

	if len(dependencyVersions) == 0 {
		return fmt.Errorf("no dependencies pending OSV scanning found")
	}

	log.Printf("Found %d dependencies pending OSV scanning\n", len(dependencyVersions))

	osvClient := NewOSVClient()
	dependencyVersionVulnerabilities, err := osvClient.GetVulnerabilityData(ctx, dependencyVersions)
	if err != nil {
		return fmt.Errorf("failed to query OSV for vulnerabilities: %w", err)
	}


	insertedVulnerabilities, err := s.repository.InsertVulnerabilities(ctx, dependencyVersionVulnerabilities)
	if err != nil {
		return fmt.Errorf("failed to insert vulnerabilities into the database: %w", err)
	}

	log.Printf("Inserted %d vulnerabilities into the database\n", len(insertedVulnerabilities))

	dependencyVersionVulnerabilityToInsert := make(map[int]int, len(insertedVulnerabilities))

	for id, vulnerabilities := range dependencyVersionVulnerabilities {
		for _, vulnerability := range vulnerabilities {
			for _, insertedVulnerability := range insertedVulnerabilities {
				if insertedVulnerability.OSVID == vulnerability.ID {
					dependencyVersionVulnerabilityToInsert[insertedVulnerability.ID] = id
					break
				}
			}
		}
	}

	if err := s.repository.InsertVersionVulnerability(ctx, dependencyVersionVulnerabilityToInsert); err != nil {
		return fmt.Errorf("failed to insert dependency version vulnerabilities: %w", err)
	}

	log.Printf("Inserted %d dependency version vulnerabilities into the database\n", len(dependencyVersionVulnerabilityToInsert))

	return nil
}
