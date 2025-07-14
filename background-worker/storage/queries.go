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
	InsertRepositoryDependencyVersionsQuery = `
		INSERT INTO repository_dependency_version (repository_id, dependency_id, version_id)
		VALUES %s
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
	GetInsertedDependenciesQuery = `
		SELECT id, name, ecosystem
		FROM dependencies
		WHERE (name, ecosystem) IN $s
	`
	GetDependenciesPendingUrlResolutionQuery = `
		SELECT D.id, D.name, D.ecosystem
		FROM dependencies as D
			INNER JOIN repository_dependency_version as RDV on D.id = RDV.dependency_id
		WHERE D.scan_status = 'pending'
			AND RDV.repository_id = $1
		ORDER BY name ASC
	`
	DeleteRepositoryDependencyVersionsQuery = `
		DELETE FROM repository_dependency_version
		WHERE repository_id = $1
	`
	InsertDependencyQuery = `
		INSERT INTO dependencies (name, ecosystem)
		VALUES %s
	`
	UpdateDependencyRepositoryIDsQuery = `
		UPDATE dependencies
		SET 
			dependency_repository_id = CASE %s END::integer,
			scanned_at = NOW(),
    	scan_status = 'completed'
		WHERE id IN (%s)
	`
	UpdateDependencyScannedFailedQuery = `
		UPDATE dependencies
		SET 
			scan_status = 'failed',
		  error_message = updates.reason
		FROM (
			SELECT unnest($1::BIGINT[]) AS id, unnest($2::TEXT[]) AS reason
		) AS updates
		WHERE dependencies.id = updates.id
	`

	// Dependency versions
	GetExistingVersionsByDependencyAndEcosystemQuery = `
		SELECT v.id, v.version, v.dependency_id, d.name as dependency_name, d.ecosystem as dependency_ecosystem
		FROM versions AS V
			INNER JOIN dependencies AS D on v.dependency_id = D.id
		WHERE
`
	InsertVersionQuery = `
		INSERT INTO versions (version, dependency_id)
		VALUES %s
	`

	// Dependency Repository
	GetDependencyRepositoryIdQuery = `
		SELECT id, repository_url FROM dependency_repository
		WHERE repository_url IN (%s);
	`
	InsertDependencyRepositoryQuery = `
		INSERT INTO dependency_repository (repository_url)
		VALUES %s
		ON CONFLICT (repository_url) DO NOTHING
	`
)
