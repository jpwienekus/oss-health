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
	db *sql.DB
}

func NewPostgresStore(connectionString string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) GetRepositoriesBySchedule(day int, hour int) ([]models.Repository, error) {
	rows, err := s.db.Query(GetRepositoriesBySchedule, day, hour)
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
	_, err := s.db.Exec(UpdateRepositoriesAsScannedQuery, id)
	if err != nil {
		return fmt.Errorf("failed to mark repository as scanned: %v", err)
	}

	return nil
}

func (s *PostgresStore) MarkRepositoryAsFailed(id int, message string) error {
	_, err := s.db.Exec(UpdateRepositoriesAsScannedFailedQuery, id, message)
	if err != nil {
		return fmt.Errorf("failed to mark repository as failed: %v", err)
	}

	return nil
}

func (s *PostgresStore) DeleteRepositoryDependencyVersion(repositoryId int) error {
	_, err := s.db.Exec(DeleteRepositoryDependencyVersionsQuery, repositoryId)
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

	rows, err := s.db.Query(query, queryArguments...)
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
	_, err := s.db.Exec(query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert dependencies: %v", err)
	}

	return nil
}

func (s *PostgresStore) GetExistingDependencyVersions(dependencies []parser.Dependency) ([]models.Dependency, error) {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("(version = $%d AND dependency_id = $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Version, d.id)
	}

	query := GetExistingDependenciesQuery
	query += strings.Join(clauses, " OR ")

	rows, err := s.db.Query(query, queryArguments...)
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
