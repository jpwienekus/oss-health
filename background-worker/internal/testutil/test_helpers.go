package testutil

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var TestCtx = context.Background()

func ClearTables(pool *pgxpool.Pool) {
	_, err := pool.Exec(TestCtx, "TRUNCATE TABLE dependencies RESTART IDENTITY CASCADE")

	if err != nil {
		log.Fatalf("failed to truncate dependencies: %v", err)
	}

	_, err = pool.Exec(TestCtx, "TRUNCATE TABLE dependency_repository RESTART IDENTITY CASCADE")

	if err != nil {
		log.Fatalf("failed to truncate dependency_repository: %v", err)
	}
}

func SeedDependencies(pool *pgxpool.Pool) {
	_, err := pool.Exec(TestCtx, `
		INSERT INTO dependencies (name, ecosystem, github_url_resolved, github_url_resolve_failed)
		VALUES
		('react', 'npm', false, false),
		('express', 'npm', false, false),
		('flask', 'pypi', false, false)
	`)

	if err != nil {
		log.Fatalf("failed to seed dependencies: %v", err)
	}
}
