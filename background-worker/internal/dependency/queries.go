package dependency

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/oss-health/background-worker/internal/db"
)

var GetPendingDependencies = func(ctx context.Context, batchSize, offset int, ecosystem string) ([]Dependency, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, ecosystem
		FROM dependencies
		WHERE github_url_resolved = false
		AND LOWER(ecosystem) = LOWER($1)
		AND github_url_resolve_failed = false
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

var UpsertGithubURLs = func(ctx context.Context, urls []string) (map[string]int64, error) {
	if len(urls) == 0 {
		return nil, nil
	}

	valueStrings := make([]string, 0, len(urls))
	valueArgs := make([]any, 0, len(urls))

	for i, url := range urls {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i+1))
		valueArgs = append(valueArgs, url)
	}

	query := fmt.Sprintf(`
		INSERT INTO dependency_repository (github_url)
		VALUES %s
		ON CONFLICT (github_url) DO NOTHING
		RETURNING id, github_url
	`, strings.Join(valueStrings, ","))

	rows, err := db.Pool.Query(ctx, query, valueArgs...)

	if err != nil {
		log.Printf("Could not insert url")
		return nil, err
	}

	defer rows.Close()

	urlToID := make(map[string]int64)

	for rows.Next() {
		var id int64
		var url string

		if err := rows.Scan(&id, &url); err != nil {
			return nil, err
		}

		urlToID[url] = id
	}

	missingURLs := []string{}

	for _, url := range urls {
		if _, ok := urlToID[url]; !ok {
			missingURLs = append(missingURLs, url)
		}
	}

	if len(missingURLs) > 0 {
		query = `SELECT id, github_url FROM dependency_repository WHERE github_url = ANY($1)`
		rows, err = db.Pool.Query(ctx, query, missingURLs)

		if err != nil {
			log.Printf("Could not read github urls: %s", err)
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var id int64
			var url string

			if err := rows.Scan(&id, &url); err != nil {
				return nil, err
			}

			urlToID[url] = id
		}
	}

	return urlToID, nil
}

var BatchUpdateDependencies = func(ctx context.Context, deps []Dependency, urlToID map[string]int64, resolvedURLs map[int64]string) error {
	tx, err := db.Pool.Begin(ctx)

	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("rollback failed: %v", err)
		}
	}()

	for _, dep := range deps {
		url, ok := resolvedURLs[dep.ID]

		if !ok {
			continue
		}

		githubURLID, ok := urlToID[url]

		if !ok {
			return fmt.Errorf("missing github_url_id for url %s", url)
		}

		if _, err := tx.Exec(ctx, `
            UPDATE dependencies
            SET dependency_repository_id = $1,
                github_url_resolved = true,
                github_url_checked_at = NOW()
            WHERE id = $2
        `, githubURLID, dep.ID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

var MarkDependenciesAsFailed = func(ctx context.Context, failureReasons map[int64]string) error {
	if len(failureReasons) == 0 {
		return nil
	}

	ids := make([]int64, 0, len(failureReasons))
	reasons := make([]string, 0, len(failureReasons))

	for id, reason := range failureReasons {
		ids = append(ids, id)
		reasons = append(reasons, reason)
	}

	_, err := db.Pool.Exec(ctx, `
		UPDATE dependencies
		SET github_url_resolve_failed = true,
		    github_url_resolve_failed_reason = updates.reason
		FROM (
			SELECT unnest($1::BIGINT[]) AS id, unnest($2::TEXT[]) AS reason
		) AS updates
		WHERE dependencies.id = updates.id
	`, ids, reasons)

	return err
}
