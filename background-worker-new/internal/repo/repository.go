package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oss-health/background-worker/internal/db"
)

type RepoRepository struct {
	db *pgxpool.Pool
}

func NewRepoRepository(db *pgxpool.Pool) *RepoRepository {
	return &RepoRepository{db: db}
}

func (r *RepoRepository) GetReposForSchedule(ctx context.Context, day int, hour int) ([]db.Repository, error) {
	rows, err := r.db.Query(ctx, GetReposForScheduleQuery, day, hour)
	if err != nil {
		return nil, fmt.Errorf("query repositories for schedule: %w", err)
	}

	defer rows.Close()

	var repositories []db.Repository

	for rows.Next() {
		var repository db.Repository
		if err := rows.Scan(&repository.ID, &repository.URL, &repository.LastScannedAt, &repository.ScanStatus); err != nil {
			return nil, fmt.Errorf("scan repository row: %w", err)
		}

		repositories = append(repositories, repository)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return repositories, nil
}

func (r *RepoRepository) MarkScanned(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, UpdateRepositoriesAsScannedQuery, id)
	if err != nil {
		return fmt.Errorf("mark repository %d as scanned: %w", id, err)
	}

	return nil
}

func (r *RepoRepository) MarkFailed(ctx context.Context, id int, message string) error {
	_, err := r.db.Exec(ctx, UpdateRepositoriesAsScannedFailedQuery, message, id)
	if err != nil {
		return fmt.Errorf("mark repository %d as failed with message %q: %w", id, message, err)
	}

	return nil
}

// TODO: continue tests from here
func (r *RepoRepository) DeleteRepositoryDependencyVersion(ctx context.Context, repositoryId int) error {
	_, err := r.db.Exec(ctx, DeleteRepositoryDependencyVersionsQuery, repositoryId)
	if err != nil {
		return fmt.Errorf("failed to delete repository dependency versions: %v", err)
	}

	return nil
}

func (r *RepoRepository) InsertRepositoryDependencyVersion(ctx context.Context, repositoryId int, dependencies []db.DependencyVersion) error {
	queryArguments := []any{}
	clauses := []string{}

	for i, d := range dependencies {
		n := i*3 + 1
		clauses = append(clauses, fmt.Sprintf("($%d, $%d, $%d)", n, n+1, n+2))
		queryArguments = append(queryArguments, repositoryId, d.DependencyId, d.ID)
	}

	query := fmt.Sprintf(InsertRepositoryDependencyVersionsQuery, strings.Join(clauses, ", "))
	_, err := r.db.Exec(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert repository dependency versions: %v", err)
	}

	return nil
}


