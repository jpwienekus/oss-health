package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/oss-health/background-worker/internal/models"
	"github.com/oss-health/background-worker/parser"

	_ "github.com/lib/pq"
)

// TODO: need pgxpool?
type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(connectionString string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	return &PostgresStore{DB: db}, nil
}

func (s *PostgresStore) GetRepositoriesBySchedule(day int, hour int) ([]models.Repository, error) {
	rows, err := s.DB.Query(GetRepositoriesBySchedule, day, hour)
	if err != nil {
		return nil, fmt.Errorf("failed to query repositories: %v", err)
	}

	defer rows.Close()
	var repositories []models.Repository

	for rows.Next() {
		var repository models.Repository
		if err := rows.Scan(&repository.ID, &repository.URL); err != nil {
			return nil, fmt.Errorf("failed to scan repository: %v", err)
		}
		repositories = append(repositories, repository)
	}

	return repositories, nil

}

func (s *PostgresStore) MarkRepositoryAsScanned(id int) error {
	_, err := s.DB.Exec(UpdateRepositoriesAsScannedQuery, id)
	if err != nil {
		return fmt.Errorf("failed to mark repository as scanned: %v", err)
	}

	return nil
}

func (s *PostgresStore) MarkRepositoryAsFailed(id int, message string) error {
	_, err := s.DB.Exec(UpdateRepositoriesAsScannedFailedQuery, message, id)
	if err != nil {
		return fmt.Errorf("failed to mark repository as failed: %v", err)
	}

	return nil
}

func (s *PostgresStore) DeleteRepositoryDependencyVersion(repositoryId int) error {
	_, err := s.DB.Exec(DeleteRepositoryDependencyVersionsQuery, repositoryId)
	if err != nil {
		return fmt.Errorf("failed to delete repository dependency versions: %v", err)
	}

	return nil
}

func (s *PostgresStore) GetExistingDependencies(dependencies []parser.Dependency) ([]models.Dependency, error) {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("(name = $%d AND ecosystem = $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Name, d.Ecosystem)
	}

	query := GetExistingDependenciesQuery
	query += strings.Join(clauses, " OR ")

	rows, err := s.DB.Query(query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing dependencies: %v", err)
	}

	defer rows.Close()

	returnValue := make([]models.Dependency, 0)

	for rows.Next() {
		var dependency models.Dependency
		if err := rows.Scan(&dependency.ID, &dependency.Name, &dependency.Ecosystem); err != nil {
			return nil, fmt.Errorf("failed to scan existing dependency: %v", err)
		}

		returnValue = append(returnValue, dependency)
	}

	return returnValue, nil
}

func (s *PostgresStore) InsertDependencies(dependencies []parser.Dependency) error {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("($%d, $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Name, d.Ecosystem)
	}

	query := fmt.Sprintf(InsertDependencyQuery, strings.Join(clauses, ", "))
	_, err := s.DB.Exec(query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert dependencies: %v", err)
	}

	return nil
}

func (s *PostgresStore) GetExistingVersionsForDependencies(dependencies []parser.Dependency) ([]models.DependencyVersion, error) {
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

	rows, err := s.DB.Query(query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing dependencies: %v", err)
	}

	defer rows.Close()

	returnValue := make([]models.DependencyVersion, 0)

	for rows.Next() {
		var dependencyVersion models.DependencyVersion
		if err := rows.Scan(&dependencyVersion.ID, &dependencyVersion.Version, &dependencyVersion.DependencyId, &dependencyVersion.DependencyName, &dependencyVersion.DependencyEcosystem); err != nil {
			return nil, fmt.Errorf("failed to scan existing dependency: %v", err)
		}

		returnValue = append(returnValue, dependencyVersion)
	}

	return returnValue, nil
}

func (s *PostgresStore) InsertDependencyVersions(dependencies []models.DependencyVersion) error {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("($%d, $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Version, d.DependencyId)
	}

	query := fmt.Sprintf(InsertVersionQuery, strings.Join(clauses, ", "))
	_, err := s.DB.Exec(query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert dependency versions: %v", err)
	}

	return nil
}

func (s *PostgresStore) InsertRepositoryDependencyVersion(repositoryId int, dependencies []models.DependencyVersion) error {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*3 + 1
		clauses = append(clauses, fmt.Sprintf("($%d, $%d, $%d)", n, n+1, n+2))
		queryArguments = append(queryArguments, repositoryId, d.DependencyId, d.ID)
	}

	query := fmt.Sprintf(InsertRepositoryDependencyVersionsQuery, strings.Join(clauses, ", "))
	_, err := s.DB.Exec(query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert repository dependency versions: %v", err)
	}

	return nil
}

func (s *PostgresStore) GetDependenciesPendingUrlResolution(repositoryId int) ([]models.Dependency, error) {
	rows, err := s.DB.Query(GetDependenciesPendingUrlResolutionQuery, repositoryId)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies pending URL resolution: %v", err)
	}

	defer rows.Close()

	var dependencies []models.Dependency
	for rows.Next() {
		var dependency models.Dependency
		if err := rows.Scan(&dependency.ID, &dependency.Name, &dependency.Ecosystem); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %v", err)
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func (s *PostgresStore) InsertDependencyRepository(dependencyUrlMap map[int]string) error {
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

	_, err := s.DB.Exec(query, queryArguments...)
	if err != nil {
		return fmt.Errorf("insert dependency_repository failed: %w", err)
	}

	return nil
}

func (s *PostgresStore) GetDependencyRepositoryIds(dependencyUrlMap map[int]string) (map[string]int, error) {
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

	rows, err := s.DB.Query(query, queryArguments...)
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

func (s *PostgresStore) UpdateDependencyRepositoryIDs(depRepoMap map[int]int) error {
	if len(depRepoMap) == 0 {
		return nil
	}

	caseClauses := make([]string, 0, len(depRepoMap))
	queryArguments := make([]any, 0, len(depRepoMap)*2)
	whereIDs := make([]int, 0, len(depRepoMap))
	argIndex := 1

	for depID, repoID := range depRepoMap {
		caseClauses = append(caseClauses, fmt.Sprintf("WHEN id = $%d THEN $%d", argIndex, argIndex+1))
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

	_, err := s.DB.Exec(query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to bulk update dependencies: %w", err)
	}

	return nil
}

func (s *PostgresStore) MarkDependenciesAsFailed(failureReasons map[int]string) error {
	if len(failureReasons) == 0 {
		return nil
	}

	ids := make([]int, 0, len(failureReasons))
	reasons := make([]string, 0, len(failureReasons))

	for id, reason := range failureReasons {
		ids = append(ids, id)
		reasons = append(reasons, reason)
	}

	_, err := s.DB.Exec(UpdateDependencyScannedFailedQuery, ids, reasons)
	if err != nil {
		return fmt.Errorf("mark dependencies failed: %w", err)
	}

	return nil
}

// func logQueryWithParams(query string, args []any) {
// 	for i, arg := range args {
// 		placeholder := fmt.Sprintf("$%d", i+1)
// 		var val string
// 		switch v := arg.(type) {
// 		case string:
// 			val = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''")) // escape single quotes
// 		default:
// 			val = fmt.Sprintf("%v", v)
// 		}
// 		query = strings.Replace(query, placeholder, val, 1)
// 	}
// 	fmt.Println("Executing Query:", query)
// }
