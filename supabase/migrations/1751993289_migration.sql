BEGIN;

-- Running upgrade 2a4b2566e2a9 -> 3dc09bffd5da

ALTER TABLE dependencies RENAME github_url_checked_at TO scanned_at;

ALTER TABLE dependencies RENAME github_url_resolve_failed_reason TO error_message;

ALTER TABLE dependency_repository RENAME github_url TO repository_url;

ALTER TABLE dependencies ADD COLUMN scan_status VARCHAR DEFAULT 'pending' NOT NULL;

DROP INDEX ix_dependencies_github_url_resolve_failed;

ALTER TABLE dependencies DROP COLUMN github_url_resolved;

ALTER TABLE dependencies DROP COLUMN github_url_resolve_failed;

UPDATE alembic_version SET version_num='3dc09bffd5da' WHERE alembic_version.version_num = '2a4b2566e2a9';

COMMIT;

