package dependencies

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

type DependencyService struct {
	repository    DependencyRepositoryInterface
	resolverChain ResolverChainInterface
}

func NewDependencyService(repository DependencyRepositoryInterface, resolverChain ResolverChainInterface) *DependencyService {
	return &DependencyService{
		repository:    repository,
		resolverChain: resolverChain,
	}
}

func (s *DependencyService) InsertDependenciesThatDoNotExist(ctx context.Context, dependencies []parsers.Dependency, workerId int) error {
	type key struct {
		Name, Ecosystem string
	}

	// These are not specific to a repository, checking what is already stored
	// in the database so no duplicates are inserted
	existingDependencies, err := s.repository.GetExistingDependencies(ctx, dependencies)
	if err != nil {
		return fmt.Errorf("failed to get existing dependencies: %w", err)
	}

	existingDependenciesMap := make(map[key]bool, len(existingDependencies))

	for _, dep := range existingDependencies {
		existingDependenciesMap[key{Name: dep.Name, Ecosystem: dep.Ecosystem}] = true
	}

	dependenciesToInsert := make([]parsers.Dependency, 0, len(dependencies))

	for _, dep := range dependencies {
		if !existingDependenciesMap[key{Name: dep.Name, Ecosystem: dep.Ecosystem}] {
			dependenciesToInsert = append(dependenciesToInsert, dep)
		}
	}

	log.Printf("[worker %d] Found %d existing dependencies, %d to insert\n", workerId, len(existingDependencies), len(dependenciesToInsert))

	if len(dependenciesToInsert) > 0 {
		if err := s.repository.InsertDependencies(ctx, dependenciesToInsert); err != nil {
			return fmt.Errorf("failed to insert dependencies: %w", err)
		}
	}

	return nil
}

func (s *DependencyService) InsertVersionsThatDoNotExist(ctx context.Context, dependencies []parsers.Dependency, workerId int) error {
	type key struct {
		Name, Ecosystem string
	}

	existingDependencies, err := s.repository.GetExistingDependencies(ctx, dependencies)
	if err != nil {
		return fmt.Errorf("failed to get existing dependencies: %w", err)
	}

	dependencyIdMap := make(map[key]int, len(existingDependencies))
	for _, dep := range existingDependencies {
		dependencyIdMap[key{Name: dep.Name, Ecosystem: dep.Ecosystem}] = dep.ID
	}

	type versionKey struct {
		Version, Name, Ecosystem string
	}

	existingVersions, err := s.repository.GetExistingVersionsForDependencies(ctx, dependencies)
	if err != nil {
		return fmt.Errorf("failed to get existing versions: %w", err)
	}

	existingVersionMap := make(map[versionKey]bool, len(existingVersions))
	for _, v := range existingVersions {
		existingVersionMap[versionKey{
			Version:   v.Version,
			Name:      v.DependencyName,
			Ecosystem: v.DependencyEcosystem,
		}] = true
	}

	versionsToInsert := []db.DependencyVersion{}

	for _, dep := range dependencies {
		vk := versionKey{Version: dep.Version, Name: dep.Name, Ecosystem: dep.Ecosystem}

		if !existingVersionMap[vk] {
			id, ok := dependencyIdMap[key{Name: dep.Name, Ecosystem: dep.Ecosystem}]

			if !ok {
				return fmt.Errorf("missing dependency ID for %s (%s)", dep.Name, dep.Ecosystem)
			}

			versionsToInsert = append(versionsToInsert, db.DependencyVersion{
				Version:      dep.Version,
				DependencyId: id,
			})
		}
	}

	log.Printf("[worker %d] Found %d existing versions, %d to insert\n", workerId, len(existingVersions), len(versionsToInsert))

	if len(versionsToInsert) > 0 {
		if err := s.repository.InsertDependencyVersions(ctx, versionsToInsert); err != nil {
			return fmt.Errorf("failed to insert dependency versions: %w", err)
		}
	}

	return nil
}

func (s *DependencyService) GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.DependencyVersion, error) {
	existingVersions, err := s.repository.GetExistingVersionsForDependencies(ctx, dependencies)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing versions: %w", err)
	}

	if len(existingVersions) == 0 {
		return nil, fmt.Errorf("no existing versions found for dependencies")
	}

	return existingVersions, nil
}

func (s *DependencyService) ResolveDependencyUrls(ctx context.Context, repositoryId int) error {
	dependencies, err := s.repository.GetDependenciesPendingUrlResolution(ctx, repositoryId)
	if err != nil {
		return fmt.Errorf("failed to get dependencies pending URL resolution: %w", err)
	}

	var (
		resolved = make(map[int]string)
		failures = make(map[int]string)
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	for _, dep := range dependencies {
		depCopy := dep
		wg.Add(1)

		go func() {
			defer wg.Done()

			url, err := s.resolverChain.Resolve(ctx, depCopy)
			log.Print(url)

			mu.Lock()
			defer mu.Unlock()
			switch {
			case err != nil:
				failures[depCopy.ID] = fmt.Sprintf("resolve: %v", err)
			case url == "":
				failures[depCopy.ID] = "empty URL"
			default:
				resolved[depCopy.ID] = url
			}
		}()
	}

	wg.Wait()

	if len(resolved) > 0 {
		if err := s.repository.InsertDependencyRepository(ctx, resolved); err != nil {
			return fmt.Errorf("failed to update dependency repository URLs: %w", err)
		}

		ids, err := s.repository.GetDependencyRepositoryIds(ctx, resolved)
		if err != nil {
			return fmt.Errorf("failed to get dependency repository IDs: %w", err)
		}

		dependencyDependencyUrlIdMap := make(map[int]int, len(ids))

		for depId, url := range resolved {
			if repoId, ok := ids[url]; ok {
				dependencyDependencyUrlIdMap[depId] = repoId
			}
		}

		if err := s.repository.UpdateDependencyRepositoryIDs(ctx, dependencyDependencyUrlIdMap); err != nil {
			return fmt.Errorf("failed to insert dependency dependency URLs: %w", err)
		}
	}

	if len(failures) > 0 {
		s.repository.MarkDependenciesAsFailed(ctx, failures)
	}

	return nil
}
