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

func (r *PostgresRepository) GetDependenciesPendingUrlResolution(ctx context.Context, batchSize, offset int, ecosystem string) ([]Dependency, error) {
	rows, err := r.db.Query(ctx, GetDependenciesPendingUrlResolutionQuery, ecosystem, offset, batchSize)

	if err != nil {
		return nil, fmt.Errorf("query pending dependencies: %w", err)
	}

	defer rows.Close()

	dependencies := make([]Dependency, 0, batchSize)

	for rows.Next() {
		var dep Dependency

		if err := rows.Scan(&dep.ID, &dep.Name, &dep.Ecosystem); err != nil {
			return nil, fmt.Errorf("scan dependency row: %w", err)
		}

		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

func (r *PostgresRepository) UpsertRepositoryURLs(ctx context.Context, resolvedUrls map[int64]string) (map[int64]int64, error) {
	if len(resolvedUrls) == 0 {
		return nil, nil
	}

	batch := &pgx.Batch{}

	for _, url := range resolvedUrls {
		batch.Queue(UpsertDependencyRepositoryQuery, url)
	}

	br := r.db.SendBatch(ctx, batch)
	defer func() {
		if err := br.Close(); err != nil {
			log.Printf("error closing reader: %v", err)
		}
	}()

	dependencyRepositoryIdUrlMap := make(map[string]int64)

	for i := 0; i < len(resolvedUrls); i++ {
		var id int64
		var url string

		if err := br.QueryRow().Scan(&id, &url); err != nil {
			return nil, fmt.Errorf("query dependency repository: %w", err)
		}

		dependencyRepositoryIdUrlMap[url] = id
	}

	dependencyDependencyRepositoryIdMap := make(map[int64]int64)

	for dependencyId, url := range resolvedUrls {
		dependencyDependencyRepositoryIdMap[dependencyId] = dependencyRepositoryIdUrlMap[url]
	}

	return dependencyDependencyRepositoryIdMap, nil
}

func (r *PostgresRepository) BatchUpdateDependencies(ctx context.Context, dependencyDependencyRepositoryIdMap map[int64]int64) error {
	batch := &pgx.Batch{}

	if len(dependencyDependencyRepositoryIdMap) == 0 {
		return nil
	}

	for dependencyId, dependencyRepositoryId := range dependencyDependencyRepositoryIdMap {
		batch.Queue(UpdateDependencyScannedQuery, dependencyRepositoryId, dependencyId)
	}

	br := r.db.SendBatch(ctx, batch)
	defer func() {
		if err := br.Close(); err != nil {
			log.Printf("error closing reader: %v", err)
		}
	}()

	for i := 0; i < len(dependencyDependencyRepositoryIdMap); i++ {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("batch update index %d: %w", i, err)
		}
	}

	return nil
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

	if _, err := r.db.Exec(ctx, UpdateDependencyScannedFailedQuery, ids, reasons); err != nil {
		return fmt.Errorf("mark dependencies failed: %w", err)
	}

	return nil
}

func (r *PostgresRepository) ReplaceRepositoryDependencyVersions(ctx context.Context, repositoryID int, pairs []DependencyVersionPair) ([]DependencyVersionResult, error) {
	var results []DependencyVersionResult

	if _, err := r.db.Exec(ctx, DeleteRepositoryDependencyVersionsQuery, repositoryID); err != nil {
		return nil, fmt.Errorf("delete repository dependency versions: %w", err)
	}

	dependenciesByKey, err := r.GetOrCreateDependencies(ctx, pairs)

	if err != nil {
		return nil, fmt.Errorf("get or create dependencies: %w", err)
	}

	versionIds, err := r.GetOrCreateVersions(ctx, pairs, dependenciesByKey)

	if err != nil {
		return nil, fmt.Errorf("get or create versions: %w", err)
	}

	var triplets [][2]int

	for _, pair := range pairs {
		depKey := pair.Name + "|" + pair.Ecosystem
		depID, ok := dependenciesByKey[depKey]

		if !ok {
			return nil, fmt.Errorf("dependency not found: %s", depKey)
		}

		versionKey := fmt.Sprintf("%d|%s", depID, pair.Version)
		verID, ok := versionIds[versionKey]

		if !ok {
			return nil, fmt.Errorf("version not found: %s", versionKey)
		}

		triplets = append(triplets, [2]int{depID, verID})
		results = append(results, DependencyVersionResult{
			VersionID: verID,
			Name:      pair.Name,
			Version:   pair.Version,
			Ecosystem: pair.Ecosystem,
		})
	}

	if err := insertRepositoryDependencyVersions(ctx, r.db, repositoryID, triplets); err != nil {
		return nil, fmt.Errorf("insert repository versions: %w", err)
	}

	return results, nil
}

func (r *PostgresRepository) GetOrCreateDependencies(ctx context.Context, pairs []DependencyVersionPair) (map[string]int, error) {
	type key struct {
		Name, Ecosystem string
	}

	seen := map[key]struct{}{}

	for _, p := range pairs {
		seen[key{p.Name, p.Ecosystem}] = struct{}{}
	}

	keys := make([]key, 0, len(seen))

	for k := range seen {
		keys = append(keys, k)
	}

	query := GetExistingDependenciesQuery
	args := []any{}
	clauses := []string{}

	for i, k := range keys {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("(name = $%d AND ecosystem = $%d)", n, n+1))
		args = append(args, k.Name, k.Ecosystem)
	}

	query += strings.Join(clauses, " OR ")

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query dependencies: %w", err)
	}

	defer rows.Close()

	result := make(map[string]int)
	existing := make(map[key]bool)

	for rows.Next() {
		var id int
		var name, ecosystem string
		if err := rows.Scan(&id, &name, &ecosystem); err != nil {
			return nil, fmt.Errorf("scan dependency: %w", err)
		}

		k := key{name, ecosystem}
		existing[k] = true
		result[name+"|"+ecosystem] = id
	}

	var insertKeys []key

	for _, k := range keys {
		if !existing[k] {
			insertKeys = append(insertKeys, k)
		}
	}

	if len(insertKeys) > 0 {
		var (
			valueStrings []string
			valueArgs    []any
		)

		for i, k := range insertKeys {
			n := i*2 + 1
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", n, n+1))
			valueArgs = append(valueArgs, k.Name, k.Ecosystem)
		}

		insertQuery := fmt.Sprintf(InsertDependencyQuery, strings.Join(valueStrings, ", "))

		rows, err := r.db.Query(ctx, insertQuery, valueArgs...)
		if err != nil {
			return nil, fmt.Errorf("insert dependencies: %w", err)
		}

		defer rows.Close()

		for rows.Next() {
			var id int
			var name, ecosystem string
			if err := rows.Scan(&id, &name, &ecosystem); err != nil {
				return nil, fmt.Errorf("scanning dependency: %w", err)
			}

			result[name+"|"+ecosystem] = id
		}
	}

	log.Printf(
		"Dependencies: %d inserted, %d existing",
		len(insertKeys), len(existing),
	)

	return result, nil
}

func (r *PostgresRepository) GetOrCreateVersions(ctx context.Context, pairs []DependencyVersionPair, depIDs map[string]int) (map[string]int, error) {
	type key struct {
		DependencyID int
		Version      string
	}

	seen := make(map[key]struct{})

	for _, p := range pairs {
		depKey := p.Name + "|" + p.Ecosystem
		depID, ok := depIDs[depKey]

		if !ok {
			return nil, fmt.Errorf("missing dependency ID: %s", depKey)
		}

		seen[key{depID, p.Version}] = struct{}{}
	}

	keys := make([]key, 0, len(seen))

	for k := range seen {
		keys = append(keys, k)
	}

	clauses := []string{}
	args := []any{}

	for i, k := range keys {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("(version = $%d AND dependency_id = $%d)", n, n+1))
		args = append(args, k.Version, k.DependencyID)
	}

	result := make(map[string]int)
	existing := make(map[key]bool)

	if len(clauses) > 0 {
		query := GetExistingVersionsQuery + strings.Join(clauses, " OR ")

		rows, err := r.db.Query(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("query existing versions: %w", err)
		}

		defer rows.Close()

		for rows.Next() {
			var id, depID int
			var ver string

			if err := rows.Scan(&id, &ver, &depID); err != nil {
				return nil, fmt.Errorf("scan version: %w", err)
			}

			k := key{depID, ver}
			existing[k] = true
			result[fmt.Sprintf("%d|%s", depID, ver)] = id
		}
	}

	var insertKeys []key

	for _, k := range keys {
		if !existing[k] {
			insertKeys = append(insertKeys, k)
		}
	}

	if len(insertKeys) > 0 {
		valueStrings := []string{}
		valueArgs := []any{}

		for i, k := range insertKeys {
			n := i*2 + 1
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", n, n+1))
			valueArgs = append(valueArgs, k.Version, k.DependencyID)
		}

		insertQuery := fmt.Sprintf(InsertVersionQuery, strings.Join(valueStrings, ", "))

		if _, err := r.db.Exec(ctx, insertQuery, valueArgs...); err != nil {
			return nil, fmt.Errorf("insert versions: %w", err)
		}

		clauses = []string{}
		args = []any{}

		for i, k := range keys {
			n := i*2 + 1
			clauses = append(clauses, fmt.Sprintf("(version = $%d AND dependency_id = $%d)", n, n+1))
			args = append(args, k.Version, k.DependencyID)
		}

		selectQuery := GetExistingVersionsQuery + strings.Join(clauses, " OR ")

		rows, err := r.db.Query(ctx, selectQuery, args...)

		if err != nil {
			return nil, fmt.Errorf("query all versions: %w", err)
		}

		defer rows.Close()

		for rows.Next() {
			var id, depID int
			var ver string

			if err := rows.Scan(&id, &ver, &depID); err != nil {
				return nil, fmt.Errorf("scan inserted version: %w", err)
			}

			result[fmt.Sprintf("%d|%s", depID, ver)] = id
		}
	}

	log.Printf(
		"Versions: %d inserted, %d existing",
		len(insertKeys), len(existing),
	)

	return result, nil
}

func insertRepositoryDependencyVersions(ctx context.Context, db *pgxpool.Pool, repositoryID int, triplets [][2]int) error {
	if len(triplets) == 0 {
		return nil
	}

	var (
		valueStrings []string
		valueArgs    []any
	)

	for i, pair := range triplets {
		n := i*3 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", n, n+1, n+2))
		valueArgs = append(valueArgs, repositoryID, pair[0], pair[1])
	}

	query := fmt.Sprintf(InsertRepositoryDependencyVersionsQuery, strings.Join(valueStrings, ", "))

	if _, err := db.Exec(ctx, query, valueArgs...); err != nil {
		return fmt.Errorf("insert repository versions: %w", err)
	}

	return nil
}
