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

	// Dependencies
	GetExistingDependenciesQuery = `
		SELECT id, name, ecosystem FROM dependencies WHERE 
	`
	DeleteRepositoryDependencyVersionsQuery = `
		DELETE FROM repository_dependency_version
		WHERE repository_id = $1
	`
	InsertDependencyQuery = `
		INSERT INTO dependencies (name, ecosystem)
		VALUES %s
	`

	// Dependency versions
	GetExistingVersionsQuery = `
		SELECT id, version, dependency_id FROM versions
		WHERE
`
)
