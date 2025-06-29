package dependency

const (
	DeleteRepositoryDependencyVersionsQuery = `
		DELETE FROM repository_dependency_version
		WHERE repository_id = $1
	`
	GetPendingDependenciesQuery = `
		SELECT id, name, ecosystem
		FROM dependencies
		WHERE github_url_resolved = false
		AND LOWER(ecosystem) = LOWER($1)
		AND github_url_resolve_failed = false
		OFFSET $2 LIMIT $3
	`
	GetDependencyIdByNameAndEcosystemQuery = `
		SELECT id FROM dependencies WHERE name = $1 AND ecosystem = $2
	`
	GetMissingUrlsQuery = `
		SELECT id, github_url 
		FROM dependency_repository 
		WHERE github_url = ANY($1)
	`
	GetSpecificVersionForDependencyQuery = `
		SELECT id FROM versions WHERE version = $1 AND dependency_id = $2
	`
	InsertDependencyQuery = `
		INSERT INTO dependencies (name, ecosystem)
		VALUES ($1, $2)
		RETURNING id
	`
	InsertDependencyRepositoryQuery = `
		INSERT INTO dependency_repository (github_url)
		VALUES %s
		ON CONFLICT (github_url) DO NOTHING
		RETURNING id, github_url
	`
	InsertRepositoryDependencyVersionsQuery = `
		INSERT INTO repository_dependency_version (repository_id, dependency_id, version_id)
		VALUES ($1, $2, $3)
	`
	InsertVersionQuery = `
		INSERT INTO versions (version, dependency_id)
		VALUES ($1, $2)
		RETURNING id
	`

	UpdateDependencyScannedQuery = `
    UPDATE dependencies
    SET dependency_repository_id = $1,
    	github_url_resolved = true,
      github_url_checked_at = NOW()
    WHERE id = $2
	`
	UpdateDependencyScannedFailedQuery = `
		UPDATE dependencies
		SET github_url_resolve_failed = true,
		    github_url_resolve_failed_reason = updates.reason
		FROM (
			SELECT unnest($1::BIGINT[]) AS id, unnest($2::TEXT[]) AS reason
		) AS updates
		WHERE dependencies.id = updates.id
	`
)
