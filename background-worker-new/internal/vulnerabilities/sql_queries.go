package vulnerabilities

const (
	InsertVulnerabilityQuery = `
		INSERT INTO vulnerabilities (osv_id, summary, severity)
		VALUES %s
		ON CONFLICT (osv_id) DO NOTHING
		RETURNING id, osv_id
	`
	InsertVersionVulnerabilityQuery = `
		INSERT INTO version_vulnerability (vulnerability_id, version_id)
		VALUES %s
	`


)
