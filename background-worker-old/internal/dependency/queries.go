package dependency

const (
	DeleteRepositoryDependencyVersionsQuery = `
		DELETE FROM repository_dependency_version
		WHERE repository_id = $1
	`
	GetDependenciesPendingUrlResolutionQuery = `
		SELECT id, name, ecosystem
		FROM dependencies
		WHERE scan_status = 'pending'
		AND LOWER(ecosystem) = LOWER($1)
		ORDER BY name ASC
		OFFSET $2 LIMIT $3
	`
	GetExistingDependenciesQuery = `
		SELECT id, name, ecosystem FROM dependencies WHERE 
	`
	GetExistingVersionsQuery = `
		SELECT id, version, dependency_id FROM versions
		WHERE
	`
	GetDependencyIdByNameAndEcosystemQuery = `
		SELECT id FROM dependencies WHERE name = $1 AND ecosystem = $2
	`
	GetSpecificVersionForDependencyQuery = `
		SELECT id FROM versions WHERE version = $1 AND dependency_id = $2
	`
	InsertDependencyQuery = `
		INSERT INTO dependencies (name, ecosystem)
		VALUES %s
		RETURNING id, name, ecosystem
	`
	InsertVersionQuery = `
		INSERT INTO versions (version, dependency_id)
		VALUES %s
	`
	InsertRepositoryDependencyVersionsQuery = `
		INSERT INTO repository_dependency_version (repository_id, dependency_id, version_id)
		VALUES %s
	`
	UpdateDependencyScannedQuery = `
    UPDATE dependencies
    SET 
			dependency_repository_id = $1,
    	scan_status = 'completed',
      scanned_at = NOW()
    WHERE id = $2
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
	UpsertDependencyRepositoryQuery = `
		INSERT INTO dependency_repository (repository_url)
		VALUES ($1)
		ON CONFLICT (repository_url) DO UPDATE SET repository_url = EXCLUDED.repository_url
		RETURNING id, repository_url
	`
)
