package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepoRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetReposForSchedule(ctx context.Context, day int, hour int) ([]Repository, error) {
	rows, err := r.db.Query(ctx, GetReposForScheduleQuery, day, hour)
	if err != nil {
		return nil, fmt.Errorf("query repositories for schedule: %w", err)
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
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return repositories, nil
}
