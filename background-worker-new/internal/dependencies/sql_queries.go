package dependencies

const (
	GetExistingDependenciesQuery = `
		SELECT id, name, ecosystem FROM dependencies WHERE 
	`
	GetDependenciesPendingUrlResolutionQuery = `
		SELECT D.id, D.name, D.ecosystem
		FROM dependencies as D
			INNER JOIN repository_dependency_version as RDV on D.id = RDV.dependency_id
		WHERE D.scan_status = 'pending'
			AND RDV.repository_id = $1
		ORDER BY name ASC
	`
	GetDependenciesPendingOSVScanQuery = `
		SELECT V.id, V.version, V.dependency_id, D.name as dependency_name, D.ecosystem as dependency_ecosystem
		FROM dependencies as D
			INNER JOIN versions as V ON D.id = V.dependency_id
		WHERE now() - COALESCE(osv_scanned_at, '1970-01-01') > interval '1 day'
	`
	InsertDependencyQuery = `
		INSERT INTO dependencies (name, ecosystem)
		VALUES %s
		ON CONFLICT (name, ecosystem) DO NOTHING
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

	// Dependency verion
	GetExistingVersionsByDependencyAndEcosystemQuery = `
		SELECT v.id, v.version, v.dependency_id, d.name as dependency_name, d.ecosystem as dependency_ecosystem
		FROM versions AS V
			INNER JOIN dependencies AS D on v.dependency_id = D.id
		WHERE
`
	InsertVersionQuery = `
		INSERT INTO versions (version, dependency_id)
		VALUES %s
		ON CONFLICT (version, dependency_id) DO NOTHING
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
	UpdateDependencyRepositoryIDsQuery = `
		UPDATE dependencies
		SET 
			dependency_repository_id = (CASE %s END),
			scanned_at = NOW(),
    	scan_status = 'completed'
		WHERE id IN (%s)
	`
)
