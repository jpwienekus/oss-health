package storage

import (
	"database/sql"
	"fmt"

	"github.com/oss-health/background-worker/internal/models"

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
