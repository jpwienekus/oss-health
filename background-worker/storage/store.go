package storage

import "github.com/oss-health/background-worker/internal/models"

type Store interface {
	GetRepositoriesBySchedule(day int, hour int) ([]models.Repository, error)
}
