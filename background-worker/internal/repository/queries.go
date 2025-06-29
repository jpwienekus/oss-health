package repository

const (
	GetRepositoriesForDayQuery = `
		SELECT id, url, last_scanned_at, scan_status
		FROM repositories
		WHERE scan_day = $1
			AND scan_hour = $2
		ORDER BY id
		LIMIT 100
	`
	UpdateRepositoriesAsScannedQuery = `
		UPDATE repositories 
		SET 
			last_scanned_at = now(), 
			scan_status = 'done'
		WHERE id = $1
	`
	UpdateRepositoriesAsScannedFailedQuery = `
		UPDATE repositories 
		SET 
			scan_status = 'error',
			error_message = $1
		WHERE id = $2
	`
)

