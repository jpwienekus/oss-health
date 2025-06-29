package repository

import (
	"context"

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
		return nil, err
	}

	defer rows.Close()

	var repositories []Repository

	for rows.Next() {
		var repository Repository
		err := rows.Scan(&repository.ID, &repository.URL, &repository.LastScannedAt, &repository.ScanStatus)

		if err != nil {
			return nil, err
		}

		repositories = append(repositories, repository)
	}

	return repositories, err
}

func (r *PostgresRepository) MarkScanned(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, UpdateRepositoriesAsScannedQuery, id)

	return err
}

func (r *PostgresRepository) MarkFailed(ctx context.Context, id int, message string) error {
	_, err := r.db.Exec(ctx, UpdateRepositoriesAsScannedFailedQuery, message, id)

	return err
}
