package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(ctx context.Context, connStr string) error {
	var err error
	Pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}

	if err := Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

func GetPendingDependencies(ctx context.Context, batchSize, offset int, ecosystem string) ([]Dependency, error) {
	rows, err := Pool.Query(ctx, `
		SELECT id, name, ecosystem
		FROM dependencies
		WHERE github_url_resolved = false
		AND LOWER(ecosystem) = LOWER($1)
		OFFSET $2 LIMIT $3
	`, ecosystem, offset, batchSize)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dependencies []Dependency

	for rows.Next() {
		var dep Dependency

		err := rows.Scan(&dep.ID, &dep.Name, &dep.Ecosystem)
		if err != nil {
			return nil, err
		}

		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

func UpdateDependencyURL(ctx context.Context, id int64, url string) error {
	_, err := Pool.Exec(ctx, `
		UPDATE dependencies
		SET github_url = $1,
		    github_url_resolved = true,
		    github_url_checked_at = NOW()
		WHERE id = $2
	`, url, id)

	return err
}
