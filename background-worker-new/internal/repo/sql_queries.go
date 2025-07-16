package repo

const (
	GetReposForScheduleQuery = `
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

	// Repository Dependency Version
	DeleteRepositoryDependencyVersionsQuery = `
		DELETE FROM repository_dependency_version
		WHERE repository_id = $1
	`
	InsertRepositoryDependencyVersionsQuery = `
		INSERT INTO repository_dependency_version (repository_id, dependency_id, version_id)
		VALUES %s
	`
)
