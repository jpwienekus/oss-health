package dependency

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

var _ DependencyRepository = (*PostgresRepository)(nil)

func (r *PostgresRepository) GetPendingDependencies(ctx context.Context, batchSize, offset int, ecosystem string) ([]Dependency, error) {
	rows, err := r.db.Query(ctx, GetPendingDependenciesQuery, ecosystem, offset, batchSize)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dependencies []Dependency

	for rows.Next() {
		var dep Dependency
		err := rows.Scan(&dep.ID, &dep.Name, &dep.Ecosystem)

		if err != nil {
			return nil, err
		}

		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

func (r *PostgresRepository) UpsertGithubURLs(ctx context.Context, urls []string) (map[string]int64, error) {
	if len(urls) == 0 {
		return nil, nil
	}

	valueStrings := make([]string, 0, len(urls))
	valueArgs := make([]any, 0, len(urls))

	for i, url := range urls {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i+1))
		valueArgs = append(valueArgs, url)
	}

	query := fmt.Sprintf(InsertDependencyRepositoryQuery, strings.Join(valueStrings, ","))

	rows, err := r.db.Query(ctx, query, valueArgs...)

	if err != nil {
		log.Printf("Could not insert url")
		return nil, err
	}

	defer rows.Close()

	urlToID := make(map[string]int64)

	for rows.Next() {
		var id int64
		var url string

		if err := rows.Scan(&id, &url); err != nil {
			return nil, err
		}

		urlToID[url] = id
	}

	missingURLs := []string{}

	for _, url := range urls {
		if _, ok := urlToID[url]; !ok {
			missingURLs = append(missingURLs, url)
		}
	}

	if len(missingURLs) > 0 {
		rows, err = r.db.Query(ctx, GetMissingUrlsQuery, missingURLs)

		if err != nil {
			log.Printf("Could not read github urls: %s", err)
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var id int64
			var url string

			if err := rows.Scan(&id, &url); err != nil {
				return nil, err
			}

			urlToID[url] = id
		}
	}

	return urlToID, nil
}

func (r *PostgresRepository) BatchUpdateDependencies(ctx context.Context, deps []Dependency, urlToID map[string]int64, resolvedURLs map[int64]string) error {
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("rollback failed: %v", err)
		}
	}()

	for _, dep := range deps {
		url, ok := resolvedURLs[dep.ID]

		if !ok {
			continue
		}

		githubURLID, ok := urlToID[url]

		if !ok {
			return fmt.Errorf("missing github_url_id for url %s", url)
		}

		if _, err := tx.Exec(ctx, UpdateDependencyScannedQuery, githubURLID, dep.ID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) MarkDependenciesAsFailed(ctx context.Context, failureReasons map[int64]string) error {
	if len(failureReasons) == 0 {
		return nil
	}

	ids := make([]int64, 0, len(failureReasons))
	reasons := make([]string, 0, len(failureReasons))

	for id, reason := range failureReasons {
		ids = append(ids, id)
		reasons = append(reasons, reason)
	}

	_, err := r.db.Exec(ctx, UpdateDependencyScannedFailedQuery, ids, reasons)

	return err
}

func (r *PostgresRepository) ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryID int, pairs []DependencyVersionPair) ([]DependencyVersionResult, error) {
	var results []DependencyVersionResult
	_, err := r.db.Exec(ctx, DeleteRepositoryDependencyVersionsQuery, repositoryID)

	if err != nil {
		return nil, fmt.Errorf("delete existing links: %w", err)
	}

	var insertedDeps, existingDeps, insertedVers, existingVers int

	for _, pair := range pairs {
		dependencyID, isNewDep, err := r.getOrCreateDependency(ctx, pair.Name, pair.Ecosystem)

		if err != nil {
			return nil, err
		}
		if isNewDep {
			insertedDeps++
		} else {
			existingDeps++
		}

		versionID, isNewVer, err := r.getOrCreateVersion(ctx, pair.Version, dependencyID)

		if err != nil {
			return nil, err
		}

		if isNewVer {
			insertedVers++
		} else {
			existingVers++
		}

		// Create the association
		_, err = r.db.Exec(ctx, InsertRepositoryDependencyVersionsQuery, repositoryID, dependencyID, versionID)
		if err != nil {
			return nil, fmt.Errorf("link repository-dependency-version: %w", err)
		}

		results = append(results, DependencyVersionResult{
			VersionID: versionID,
			Name:      pair.Name,
			Version:   pair.Version,
			Ecosystem: pair.Ecosystem,
		})
	}

	log.Printf(
		"Dependencies: %d inserted, %d existing; Versions: %d inserted, %d existing",
		insertedDeps, existingDeps, insertedVers, existingVers,
	)

	return results, nil
}

func (r *PostgresRepository) getOrCreateDependency(ctx context.Context, name, ecosystem string) (int, bool, error) {
	var id int
	err := r.db.QueryRow(ctx, GetDependencyIdByNameAndEcosystemQuery, name, ecosystem).Scan(&id)

	if err == pgx.ErrNoRows {
		err = r.db.QueryRow(ctx, InsertDependencyQuery, name, ecosystem).Scan(&id)

		if err != nil {
			return 0, false, fmt.Errorf("insert dependency %q (%s): %w", name, ecosystem, err)
		}

		return id, true, nil
	} else if err != nil {
		return 0, false, fmt.Errorf("query dependency %q (%s): %w", name, ecosystem, err)
	}

	return id, false, nil
}

func (r *PostgresRepository) getOrCreateVersion(ctx context.Context, version string, dependencyID int) (int, bool, error) {
	var id int
	err := r.db.QueryRow(ctx, GetSpecificVersionForDependencyQuery, version, dependencyID).Scan(&id)

	if err == pgx.ErrNoRows {
		err = r.db.QueryRow(ctx, InsertVersionQuery, version, dependencyID).Scan(&id)

		if err != nil {
			return 0, false, fmt.Errorf("insert version %q: %w", version, err)
		}

		return id, true, nil
	} else if err != nil {
		return 0, false, fmt.Errorf("query version %q: %w", version, err)
	}

	return id, false, nil
}
