package repository

import (
	"context"

	"github.com/oss-health/background-worker/internal/db"
)

func GetRepositoriesForDay(ctx context.Context, day int, hour int) ([]Repository, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, url, last_scanned_at, scan_status
		FROM repositories
		WHERE scan_day = $1
		AND scan_hour = $2
		ORDER BY id
		LIMIT 100
	`, day, hour)

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

var MarkScanned = func(ctx context.Context, id int) {
	_, _ = db.Pool.Exec(ctx, `
		UPDATE repositories 
		SET 
			last_scanned_at = now(), 
			scan_status = 'done'
		WHERE id = $1
		`, id)
}


var MarkFailed = func(ctx context.Context, id int, message string) {
	_, _ = db.Pool.Exec(ctx, `
		UPDATE repositories 
		SET 
			scan_status = 'error',
			error_message = $1
		WHERE id = $2
		`, message, id)
}
