package dependency

const (
	GetPendingDependenciesQuery = `
		SELECT id, name, ecosystem
		FROM dependencies
		WHERE github_url_resolved = false
		AND LOWER(ecosystem) = LOWER($1)
		AND github_url_resolve_failed = false
		OFFSET $2 LIMIT $3
	`
	GetMissingUrlsQuery = `
		SELECT id, github_url 
		FROM dependency_repository 
		WHERE github_url = ANY($1)
	`
	InsertDependencyRepositoryQuery = `
		INSERT INTO dependency_repository (github_url)
		VALUES %s
		ON CONFLICT (github_url) DO NOTHING
		RETURNING id, github_url
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
