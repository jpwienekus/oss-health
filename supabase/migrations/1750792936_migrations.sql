BEGIN;

CREATE TABLE alembic_version (
    version_num VARCHAR(32) NOT NULL, 
    CONSTRAINT alembic_version_pkc PRIMARY KEY (version_num)
);

-- Running upgrade  -> 77118a438b4e

CREATE TABLE "user" (
    id SERIAL NOT NULL, 
    github_id INTEGER NOT NULL, 
    github_username VARCHAR NOT NULL, 
    PRIMARY KEY (id), 
    UNIQUE (github_id)
);

CREATE INDEX ix_user_id ON "user" (id);

INSERT INTO alembic_version (version_num) VALUES ('77118a438b4e') RETURNING alembic_version.version_num;

-- Running upgrade 77118a438b4e -> 105d2b53c498

CREATE TABLE repositories (
    id SERIAL NOT NULL, 
    github_id INTEGER NOT NULL, 
    name VARCHAR NOT NULL, 
    description VARCHAR, 
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL, 
    user_id INTEGER NOT NULL, 
    PRIMARY KEY (id), 
    FOREIGN KEY(user_id) REFERENCES "user" (id), 
    UNIQUE (github_id)
);

CREATE INDEX ix_repositories_id ON repositories (id);

UPDATE alembic_version SET version_num='105d2b53c498' WHERE alembic_version.version_num = '77118a438b4e';

-- Running upgrade 105d2b53c498 -> 83222855fffd

ALTER TABLE "user" ADD COLUMN synced_at TIMESTAMP WITHOUT TIME ZONE;

UPDATE alembic_version SET version_num='83222855fffd' WHERE alembic_version.version_num = '105d2b53c498';

-- Running upgrade 83222855fffd -> d03f1b6980cd

ALTER TABLE repositories ADD COLUMN url VARCHAR;

ALTER TABLE repositories ADD COLUMN open_issues INTEGER;

ALTER TABLE repositories ADD COLUMN score INTEGER;

UPDATE alembic_version SET version_num='d03f1b6980cd' WHERE alembic_version.version_num = '83222855fffd';

-- Running upgrade d03f1b6980cd -> 7095f5c0ef18

ALTER TABLE "user" ADD COLUMN access_token VARCHAR;

UPDATE alembic_version SET version_num='7095f5c0ef18' WHERE alembic_version.version_num = 'd03f1b6980cd';

-- Running upgrade 7095f5c0ef18 -> 19ea01e48843

ALTER TABLE repositories DROP COLUMN url;

ALTER TABLE repositories DROP COLUMN open_issues;

ALTER TABLE repositories DROP COLUMN updated_at;

ALTER TABLE repositories DROP COLUMN name;

ALTER TABLE repositories DROP COLUMN description;

UPDATE alembic_version SET version_num='19ea01e48843' WHERE alembic_version.version_num = '7095f5c0ef18';

-- Running upgrade 19ea01e48843 -> 8bae674e01e4

ALTER TABLE "user" DROP COLUMN synced_at;

UPDATE alembic_version SET version_num='8bae674e01e4' WHERE alembic_version.version_num = '19ea01e48843';

-- Running upgrade 8bae674e01e4 -> c5864af6233f

ALTER TABLE repositories ADD COLUMN clone_url VARCHAR;

UPDATE alembic_version SET version_num='c5864af6233f' WHERE alembic_version.version_num = '8bae674e01e4';

-- Running upgrade c5864af6233f -> d1f018fd6b64

CREATE TABLE dependencies (
    id SERIAL NOT NULL, 
    name VARCHAR, 
    version VARCHAR, 
    ecosystem VARCHAR, 
    PRIMARY KEY (id)
);

CREATE INDEX ix_dependencies_id ON dependencies (id);

CREATE TABLE repository_dependency (
    repository_id INTEGER NOT NULL, 
    dependency_id INTEGER NOT NULL, 
    PRIMARY KEY (repository_id, dependency_id), 
    FOREIGN KEY(dependency_id) REFERENCES dependencies (id), 
    FOREIGN KEY(repository_id) REFERENCES repositories (id)
);

UPDATE alembic_version SET version_num='d1f018fd6b64' WHERE alembic_version.version_num = 'c5864af6233f';

-- Running upgrade d1f018fd6b64 -> 6bf0ff096894

CREATE TABLE vulnerabilities (
    id SERIAL NOT NULL, 
    osv_id VARCHAR NOT NULL, 
    PRIMARY KEY (id), 
    UNIQUE (osv_id)
);

CREATE TABLE dependency_vulnerability (
    dependency_id INTEGER NOT NULL, 
    vulnerability_id INTEGER NOT NULL, 
    PRIMARY KEY (dependency_id, vulnerability_id), 
    FOREIGN KEY(dependency_id) REFERENCES dependencies (id), 
    FOREIGN KEY(vulnerability_id) REFERENCES vulnerabilities (id)
);

UPDATE alembic_version SET version_num='6bf0ff096894' WHERE alembic_version.version_num = 'd1f018fd6b64';

-- Running upgrade 6bf0ff096894 -> c362a2813ad2

CREATE TABLE versions (
    id SERIAL NOT NULL, 
    version VARCHAR NOT NULL, 
    PRIMARY KEY (id)
);

CREATE TABLE dependency_version (
    dependency_id INTEGER NOT NULL, 
    version_id INTEGER NOT NULL, 
    PRIMARY KEY (dependency_id, version_id), 
    FOREIGN KEY(dependency_id) REFERENCES dependencies (id), 
    FOREIGN KEY(version_id) REFERENCES versions (id)
);

CREATE TABLE version_vulnerability (
    dependency_id INTEGER NOT NULL, 
    vulnerability_id INTEGER NOT NULL, 
    PRIMARY KEY (dependency_id, vulnerability_id), 
    FOREIGN KEY(dependency_id) REFERENCES versions (id), 
    FOREIGN KEY(vulnerability_id) REFERENCES vulnerabilities (id)
);

DROP TABLE dependency_vulnerability;

ALTER TABLE dependencies DROP COLUMN version;

UPDATE alembic_version SET version_num='c362a2813ad2' WHERE alembic_version.version_num = '6bf0ff096894';

-- Running upgrade c362a2813ad2 -> 0b1adaba6815

CREATE TABLE repository_dependency_version (
    repository_id INTEGER NOT NULL, 
    dependency_id INTEGER NOT NULL, 
    version_id INTEGER NOT NULL, 
    PRIMARY KEY (repository_id, dependency_id, version_id), 
    FOREIGN KEY(dependency_id) REFERENCES dependencies (id), 
    FOREIGN KEY(repository_id) REFERENCES repositories (id), 
    FOREIGN KEY(version_id) REFERENCES versions (id)
);

UPDATE alembic_version SET version_num='0b1adaba6815' WHERE alembic_version.version_num = 'c362a2813ad2';

-- Running upgrade 0b1adaba6815 -> 4a4d201d5e00

DROP TABLE dependency_version;

UPDATE alembic_version SET version_num='4a4d201d5e00' WHERE alembic_version.version_num = '0b1adaba6815';

-- Running upgrade 4a4d201d5e00 -> 80f78b644553

DROP TABLE repository_dependency;

UPDATE alembic_version SET version_num='80f78b644553' WHERE alembic_version.version_num = '4a4d201d5e00';

-- Running upgrade 80f78b644553 -> 5b4e6d22b0ac

ALTER TABLE versions ADD COLUMN dependency_id INTEGER NOT NULL;

ALTER TABLE versions ADD UNIQUE (version, dependency_id);

ALTER TABLE versions ADD FOREIGN KEY(dependency_id) REFERENCES dependencies (id);

UPDATE alembic_version SET version_num='5b4e6d22b0ac' WHERE alembic_version.version_num = '80f78b644553';

-- Running upgrade 5b4e6d22b0ac -> 142e721d7148

ALTER TABLE version_vulnerability ADD COLUMN version_id INTEGER NOT NULL;

ALTER TABLE version_vulnerability DROP CONSTRAINT version_vulnerability_dependency_id_fkey;

ALTER TABLE version_vulnerability ADD FOREIGN KEY(version_id) REFERENCES versions (id);

ALTER TABLE version_vulnerability DROP COLUMN dependency_id;

UPDATE alembic_version SET version_num='142e721d7148' WHERE alembic_version.version_num = '5b4e6d22b0ac';

-- Running upgrade 142e721d7148 -> 83d4ffb8473c

ALTER TABLE repositories ADD COLUMN scanned_date TIMESTAMP WITH TIME ZONE;

UPDATE alembic_version SET version_num='83d4ffb8473c' WHERE alembic_version.version_num = '142e721d7148';

-- Running upgrade 83d4ffb8473c -> 49823a158fb2

ALTER TABLE dependencies ADD COLUMN github_url VARCHAR;

ALTER TABLE dependencies ADD COLUMN github_url_resolved BOOLEAN DEFAULT false NOT NULL;

ALTER TABLE dependencies ADD COLUMN github_url_checked_at TIMESTAMP WITH TIME ZONE;

UPDATE alembic_version SET version_num='49823a158fb2' WHERE alembic_version.version_num = '83d4ffb8473c';

-- Running upgrade 49823a158fb2 -> 77abd35aa613

CREATE TABLE dependency_repository (
    id SERIAL NOT NULL, 
    github_url VARCHAR NOT NULL, 
    PRIMARY KEY (id), 
    UNIQUE (github_url)
);

CREATE INDEX ix_dependency_repository_id ON dependency_repository (id);

ALTER TABLE dependencies ADD COLUMN dependency_repository_id INTEGER;

ALTER TABLE dependencies ADD FOREIGN KEY(dependency_repository_id) REFERENCES dependency_repository (id);

ALTER TABLE dependencies DROP COLUMN github_url;

UPDATE alembic_version SET version_num='77abd35aa613' WHERE alembic_version.version_num = '49823a158fb2';

-- Running upgrade 77abd35aa613 -> 1f0271d16689

ALTER TABLE dependencies ADD COLUMN github_url_resolve_failed BOOLEAN DEFAULT false NOT NULL;

ALTER TABLE dependencies ADD COLUMN github_url_resolve_failed_reason VARCHAR;

CREATE INDEX ix_dependencies_github_url_resolve_failed ON dependencies (github_url_resolve_failed);

UPDATE alembic_version SET version_num='1f0271d16689' WHERE alembic_version.version_num = '77abd35aa613';

COMMIT;

