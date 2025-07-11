package storage

const (
	// Repository queries
	GetRepositoriesBySchedule = `
		SELECT id, url
		FROM repositories
		WHERE scan_day = $1
			AND scan_hour = $2
		ORDER BY id
	`
)
