package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepositoryRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

var _ RepositoryRepository = (*PostgresRepository)(nil)

func (r *PostgresRepository) GetRepositoriesForDay(ctx context.Context, day int, hour int) ([]Repository, error) {
	rows, err := r.db.Query(ctx, GetRepositoriesForDayQuery, day, hour)
	if err != nil {
		return nil, fmt.Errorf("query repositories for day %d hour %d: %w", day, hour, err)
	}

	defer rows.Close()

	var repositories []Repository

	for rows.Next() {
		var repository Repository
		if err := rows.Scan(&repository.ID, &repository.URL, &repository.LastScannedAt, &repository.ScanStatus); err != nil {
			return nil, fmt.Errorf("scan repository row: %w", err)
		}

		repositories = append(repositories, repository)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate repository rows: %w", err)
	}

	return repositories, nil
}

func (r *PostgresRepository) MarkScanned(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, UpdateRepositoriesAsScannedQuery, id)
	if err != nil {
		return fmt.Errorf("mark repository %d as scanned: %w", id, err)
	}

	return nil
}

func (r *PostgresRepository) MarkFailed(ctx context.Context, id int, message string) error {
	_, err := r.db.Exec(ctx, UpdateRepositoriesAsScannedFailedQuery, message, id)
	if err != nil {
		return fmt.Errorf("mark repository %d as failed with message %q: %w", id, message, err)
	}

	return nil
}
