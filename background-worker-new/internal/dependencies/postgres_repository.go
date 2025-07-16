package dependencies

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oss-health/background-worker/internal/repo/parsers"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewDependencyRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetExistingDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]Dependency, error) {
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

	returnValue := make([]Dependency, 0)

	for rows.Next() {
		var dependency Dependency
		if err := rows.Scan(&dependency.ID, &dependency.Name, &dependency.Ecosystem); err != nil {
			return nil, fmt.Errorf("failed to scan existing dependency: %v", err)
		}

		returnValue = append(returnValue, dependency)
	}

	return returnValue, nil
}


func (r *PostgresRepository) InsertDependencies(ctx context.Context,  dependencies []parsers.Dependency) error {
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


func (r *PostgresRepository) GetExistingVersionsForDependencies(ctx context.Context, dependencies []parsers.Dependency) ([]DependencyVersion, error) {
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

	rows, err := r.db.Query(ctx, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing dependencies: %v", err)
	}

	defer rows.Close()

	returnValue := make([]DependencyVersion, 0)

	for rows.Next() {
		var dependencyVersion DependencyVersion
		if err := rows.Scan(&dependencyVersion.ID, &dependencyVersion.Version, &dependencyVersion.DependencyId, &dependencyVersion.DependencyName, &dependencyVersion.DependencyEcosystem); err != nil {
			return nil, fmt.Errorf("failed to scan existing dependency: %v", err)
		}

		returnValue = append(returnValue, dependencyVersion)
	}

	return returnValue, nil
}


func (r *PostgresRepository) InsertDependencyVersions(ctx context.Context, dependencies []DependencyVersion) error {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*2 + 1
		clauses = append(clauses, fmt.Sprintf("($%d, $%d)", n, n+1))
		queryArguments = append(queryArguments, d.Version, d.DependencyId)
	}

	query := fmt.Sprintf(InsertVersionQuery, strings.Join(clauses, ", "))
	_, err := r.db.Exec(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert dependency versions: %v", err)
	}

	return nil
}
