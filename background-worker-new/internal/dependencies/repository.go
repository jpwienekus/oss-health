package dependencies

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repo/parsers"
)

type DependencyRepository struct {
	db *pgxpool.Pool
}

func NewDependencyRepository(db *pgxpool.Pool) *DependencyRepository {
	return &DependencyRepository{db: db}
}

func (r *DependencyRepository) GetExistingDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.Dependency, error) {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("(name = $%d AND ecosystem = $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Name, d.Ecosystem)
	}

	query := GetExistingDependenciesQuery
	query += strings.Join(clauses, " OR ")

	rows, err := r.db.Query(ctx, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing dependencies: %v", err)
	}

	defer rows.Close()

	returnValue := make([]db.Dependency, 0)

	for rows.Next() {
		var dependency db.Dependency
		if err := rows.Scan(&dependency.ID, &dependency.Name, &dependency.Ecosystem); err != nil {
			return nil, fmt.Errorf("failed to scan existing dependency: %v", err)
		}

		returnValue = append(returnValue, dependency)
	}

	return returnValue, nil
}

func (r *DependencyRepository) InsertDependencies(ctx context.Context, dependencies []parsers.Dependency) error {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("($%d, $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Name, d.Ecosystem)
	}

	query := fmt.Sprintf(InsertDependencyQuery, strings.Join(clauses, ", "))
	_, err := r.db.Exec(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert dependencies: %v", err)
	}

	return nil
}

func (r *DependencyRepository) GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]db.DependencyVersion, error) {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*3 + 1
		// NOTE: All parsers need to call ToLower on these values,
		// Don't want to have to call LOWER in sql for each row for
		// for performance reasons.
		clauses = append(clauses, fmt.Sprintf("(v.version = $%d AND d.name = $%d AND d.ecosystem = $%d)", n, n+1, n+2))
		queryArguments = append(queryArguments, strings.ToLower(d.Version), strings.ToLower(d.Name), strings.ToLower(d.Ecosystem))
	}

	query := GetExistingVersionsByDependencyAndEcosystemQuery
	query += strings.Join(clauses, " OR ")
	logQueryWithParams(query, queryArguments)

	rows, err := r.db.Query(ctx, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing dependencies: %v", err)
	}

	defer rows.Close()

	returnValue := make([]db.DependencyVersion, 0)

	for rows.Next() {
		var dependencyVersion db.DependencyVersion
		if err := rows.Scan(&dependencyVersion.ID, &dependencyVersion.Version, &dependencyVersion.DependencyId, &dependencyVersion.DependencyName, &dependencyVersion.DependencyEcosystem); err != nil {
			return nil, fmt.Errorf("failed to scan existing dependency: %v", err)
		}

		returnValue = append(returnValue, dependencyVersion)
	}

	return returnValue, nil
}

func (r *DependencyRepository) InsertDependencyVersions(ctx context.Context, dependencies []db.DependencyVersion) error {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("($%d, $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Version, d.DependencyId)
	}

	query := fmt.Sprintf(InsertVersionQuery, strings.Join(clauses, ", "))
	logQueryWithParams(query, queryArguments)
	_, err := r.db.Exec(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert dependency versions: %v", err)
	}

	return nil
}

func (r *DependencyRepository) GetDependenciesPendingUrlResolution(ctx context.Context, repositoryId int) ([]db.Dependency, error) {
	rows, err := r.db.Query(ctx, GetDependenciesPendingUrlResolutionQuery, repositoryId)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies pending URL resolution: %v", err)
	}

	defer rows.Close()

	var dependencies []db.Dependency
	for rows.Next() {
		var dependency db.Dependency
		if err := rows.Scan(&dependency.ID, &dependency.Name, &dependency.Ecosystem); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %v", err)
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func (r *DependencyRepository) InsertDependencyRepository(ctx context.Context, dependencyUrlMap map[int]string) error {
	if len(dependencyUrlMap) == 0 {
		return nil // nothing to update
	}

	valueStrings := make([]string, len(dependencyUrlMap))
	queryArguments := make([]any, len(dependencyUrlMap))
	i := 0

	for _, url := range dependencyUrlMap {
		valueStrings[i] = fmt.Sprintf("($%d)", i+1)
		queryArguments[i] = url
		i++
	}

	query := fmt.Sprintf(InsertDependencyRepositoryQuery, strings.Join(valueStrings, ", "))

	_, err := r.db.Exec(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("insert dependency_repository failed: %w", err)
	}

	return nil
}

func (r *DependencyRepository) GetDependencyRepositoryIds(ctx context.Context, dependencyUrlMap map[int]string) (map[string]int, error) {
	if len(dependencyUrlMap) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(dependencyUrlMap))
	queryArguments := make([]any, len(dependencyUrlMap))
	i := 0

	for _, url := range dependencyUrlMap {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		queryArguments[i] = url
		i++
	}

	query := fmt.Sprintf(GetDependencyRepositoryIdQuery, strings.Join(placeholders, ", "))

	rows, err := r.db.Query(ctx, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("query dependency_repository failed: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var id int
		var url string
		if err := rows.Scan(&id, &url); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		result[url] = id
	}
	return result, nil
}

func (r *DependencyRepository) UpdateDependencyRepositoryIDs(ctx context.Context, depRepoMap map[int]int) error {
	if len(depRepoMap) == 0 {
		return nil
	}

	caseClauses := make([]string, 0, len(depRepoMap))
	queryArguments := make([]any, 0, len(depRepoMap)*2)
	whereIDs := make([]int, 0, len(depRepoMap))
	argIndex := 1

	for depID, repoID := range depRepoMap {
		caseClauses = append(caseClauses, fmt.Sprintf("WHEN id = $%d THEN $%d::integer", argIndex, argIndex+1))
		queryArguments = append(queryArguments, depID, repoID)
		whereIDs = append(whereIDs, depID)
		argIndex += 2
	}

	wherePlaceholders := make([]string, len(whereIDs))
	for i := range whereIDs {
		wherePlaceholders[i] = fmt.Sprintf("$%d", argIndex+i)
		queryArguments = append(queryArguments, whereIDs[i])
	}

	query := fmt.Sprintf(UpdateDependencyRepositoryIDsQuery, strings.Join(caseClauses, " "), strings.Join(wherePlaceholders, ", "))

	_, err := r.db.Exec(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to bulk update dependencies: %w", err)
	}

	return nil
}

func (r *DependencyRepository) MarkDependenciesAsFailed(ctx context.Context, failureReasons map[int]string) error {
	if len(failureReasons) == 0 {
		return nil
	}

	ids := make([]int, 0, len(failureReasons))
	reasons := make([]string, 0, len(failureReasons))

	for id, reason := range failureReasons {
		ids = append(ids, id)
		reasons = append(reasons, reason)
	}

	_, err := r.db.Exec(ctx, UpdateDependencyScannedFailedQuery, ids, reasons)
	if err != nil {
		return fmt.Errorf("mark dependencies failed: %w", err)
	}

	return nil
}

// TODO: continue tests from here:
func (r *DependencyRepository) GetDependenciesPendingOSVScan(ctx context.Context) ([]db.DependencyVersion, error) {
	rows, err := r.db.Query(ctx, GetDependenciesPendingOSVScanQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies pending OSV scan: %v", err)
	}

	defer rows.Close()

	var dependencyVersions []db.DependencyVersion
	for rows.Next() {
		var dependencyVersion db.DependencyVersion
		if err := rows.Scan(&dependencyVersion.ID, &dependencyVersion.Version, &dependencyVersion.DependencyId, &dependencyVersion.DependencyName, &dependencyVersion.DependencyEcosystem); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %v", err)
		}
		dependencyVersions = append(dependencyVersions, dependencyVersion)
	}


	return dependencyVersions, nil
}

func logQueryWithParams(query string, args []any) {
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		var val string
		switch v := arg.(type) {
		case string:
			val = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''")) // escape single quotes
		default:
			val = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, placeholder, val, 1)
	}
	fmt.Println("Executing Query:", query)
}
