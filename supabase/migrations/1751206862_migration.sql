BEGIN;

-- Running upgrade 77abd35aa613 -> 1f0271d16689

ALTER TABLE dependencies ADD COLUMN github_url_resolve_failed BOOLEAN DEFAULT false NOT NULL;

ALTER TABLE dependencies ADD COLUMN github_url_resolve_failed_reason VARCHAR;

CREATE INDEX ix_dependencies_github_url_resolve_failed ON dependencies (github_url_resolve_failed);

UPDATE alembic_version SET version_num='1f0271d16689' WHERE alembic_version.version_num = '77abd35aa613';

-- Running upgrade 1f0271d16689 -> 2a4b2566e2a9

ALTER TABLE repositories ADD COLUMN scan_day INTEGER DEFAULT 0 NOT NULL;

ALTER TABLE repositories ADD COLUMN scan_hour INTEGER DEFAULT 0 NOT NULL;

ALTER TABLE repositories ADD COLUMN scan_status VARCHAR DEFAULT 'pending' NOT NULL;

ALTER TABLE repositories ADD COLUMN error_message VARCHAR;

UPDATE repositories SET clone_url = '' where clone_url IS NULL;

ALTER TABLE repositories ALTER COLUMN clone_url SET NOT NULL;

ALTER TABLE repositories RENAME clone_url TO url;

ALTER TABLE repositories RENAME scanned_date TO last_scanned_at;

UPDATE alembic_version SET version_num='2a4b2566e2a9' WHERE alembic_version.version_num = '1f0271d16689';

COMMIT;

