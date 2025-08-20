package repository

import (
	"context"
)

type RepositoryRepository interface {
	GetRepositoriesForDay(ctx context.Context, day int, hour int) ([]Repository, error)
	MarkScanned(ctx context.Context, id int) error
	MarkFailed(ctx context.Context, id int, message string) error
}
