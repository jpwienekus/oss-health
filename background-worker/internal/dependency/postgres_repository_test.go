package dependency_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/stretchr/testify/assert"
)

var TestDB *pgxpool.Pool
var TestCtx = context.Background()

func ClearTables(pool *pgxpool.Pool) {
	tables := []string{
		"repository_dependency_version",
		"dependency_repository",
		"versions",
		"dependencies",
		"repositories",
		`"user"`,
	}

	for _, table := range tables {
		_, err := pool.Exec(TestCtx, fmt.Sprintf(`TRUNCATE TABLE %s CASCADE`, table))

		if err != nil {
			log.Fatalf("failed to truncate %s: %v", table, err)
		}
	}
}

func SeedDependencies(pool *pgxpool.Pool) {
	_, err := pool.Exec(TestCtx, `
		INSERT INTO dependencies (id, name, ecosystem, scan_status)
		VALUES
			(1, 'react', 'npm', 'pending'),
			(2, 'express', 'npm', 'pending'),
			(3, 'flask', 'pypi', 'pending')
	`)

	if err != nil {
		log.Fatalf("failed to seed dependencies: %v", err)
	}
}

func SeedUsers(pool *pgxpool.Pool) {
	_, err := pool.Exec(TestCtx, `
		INSERT INTO "user" (id, github_id, github_username, access_token)
		VALUES 
			(1, 101, 'user101', 'token1')
	`)

	if err != nil {
		log.Fatalf("failed to seed dependencies: %v", err)
	}
}

func SeedRepositories(pool *pgxpool.Pool) {
	_, err := pool.Exec(TestCtx, `
		INSERT INTO repositories (url, github_id, last_scanned_at, scan_status)
		VALUES ('https://github.com/test/repo', 123456, now(), 'pending')
		RETURNING id
	`)

	if err != nil {
		log.Fatalf("failed to seed dependencies: %v", err)
	}
}

func TestMain(m *testing.M) {
	connStr := "postgres://test-user:password@localhost:5434/test_db"
	var err error
	TestDB, err = db.Connect(TestCtx, connStr)

	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}

	ClearTables(TestDB)

	code := m.Run()
	TestDB.Close()
	os.Exit(code)
}

func TestGetDependenciesPendingUrlResolution(t *testing.T) {
	ClearTables(TestDB)
	SeedDependencies(TestDB)

	repository := dependency.NewPostgresRepository(TestDB)
	dependencies, err := repository.GetDependenciesPendingUrlResolution(TestCtx, 10, 0, "npm")
	assert.NoError(t, err)
	assert.Len(t, dependencies, 2)

	names := []string{dependencies[0].Name, dependencies[1].Name}
	assert.Contains(t, names, "react")
	assert.Contains(t, names, "express")
}

func TestUpsertRepositoryURLs(t *testing.T) {
	ClearTables(TestDB)

	urlMap := map[int64]string{
		1: "https://github.com/facebook/react",
		2: "https://github.com/expressjs/express",
	}

	repository := dependency.NewPostgresRepository(TestDB)
	urlToID, err := repository.UpsertRepositoryURLs(TestCtx, urlMap)
	assert.NoError(t, err)
	assert.Len(t, urlToID, 2)

	for repDepId, url := range urlMap {
		id, ok := urlToID[repDepId]
		assert.True(t, ok, fmt.Sprintf("url %s not found in result", url))
		assert.Greater(t, id, int64(0))
	}

	// Insert duplicates again, expect same IDs returned (no duplicates)
	urlToID2, err := repository.UpsertRepositoryURLs(TestCtx, urlMap)
	assert.NoError(t, err)
	assert.Equal(t, urlToID, urlToID2)
}

func TestBatchUpdateDependencies(t *testing.T) {
	ClearTables(TestDB)
	SeedDependencies(TestDB)

	urlMap := map[int64]string{
		1: "https://github.com/facebook/react",
	}

	repository := dependency.NewPostgresRepository(TestDB)
	urlToID, err := repository.UpsertRepositoryURLs(TestCtx, urlMap)
	assert.NoError(t, err)

	err = repository.BatchUpdateDependencies(TestCtx, urlToID)
	assert.NoError(t, err)

	var status string
	var id int
	err = TestDB.QueryRow(TestCtx, `SELECT id, scan_status FROM dependencies WHERE name='react'`).Scan(&id, &status)

	assert.NoError(t, err)
	assert.NotEmpty(t, status)
}

func TestMarkDependenciesAsFailed(t *testing.T) {
	ClearTables(TestDB)
	SeedDependencies(TestDB)

	repository := dependency.NewPostgresRepository(TestDB)
	deps, err := repository.GetDependenciesPendingUrlResolution(TestCtx, 10, 0, "npm")
	assert.NoError(t, err)

	failureReasons := map[int64]string{}

	for _, d := range deps {
		failureReasons[d.ID] = "Failed to resolve URL"
	}

	err = repository.MarkDependenciesAsFailed(TestCtx, failureReasons)
	assert.NoError(t, err)

	var status string
	var reason string
	err = TestDB.QueryRow(TestCtx, `SELECT scan_status, error_message FROM dependencies WHERE id=$1`, deps[0].ID).Scan(&status, &reason)

	assert.NoError(t, err)
	assert.Equal(t, "failed", status)
	assert.Equal(t, "Failed to resolve URL", reason)
}

func TestReplaceRepositoryDependencyVersions(t *testing.T) {
	ClearTables(TestDB)
	SeedUsers(TestDB)

	var repositoryID int
	err := TestDB.QueryRow(TestCtx, `
		INSERT INTO repositories (url, github_id, user_id, last_scanned_at, scan_status)
		VALUES ('https://github.com/test/repo', 123456, 1, now(), 'pending')
		RETURNING id
	`).Scan(&repositoryID)
	assert.NoError(t, err)

	repo := dependency.NewPostgresRepository(TestDB)

	pairs := []dependency.DependencyVersionPair{
		{
			Name:      "express",
			Version:   "4.18.2",
			Ecosystem: "npm",
		},
		{
			Name:      "flask",
			Version:   "2.1.0",
			Ecosystem: "pypi",
		},
	}

	results, err := repo.ReplaceRepositoryDependencyVersions(TestCtx, repositoryID, pairs)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	rows, err := TestDB.Query(TestCtx, `
		SELECT d.name, v.version, d.ecosystem
		FROM repository_dependency_version rdv
		JOIN dependencies d ON d.id = rdv.dependency_id
		JOIN versions v ON v.id = rdv.version_id
		WHERE rdv.repository_id = $1
	`, repositoryID)
	assert.NoError(t, err)

	var found []dependency.DependencyVersionPair
	for rows.Next() {
		var name, version, ecosystem string
		err := rows.Scan(&name, &version, &ecosystem)
		assert.NoError(t, err)
		found = append(found, dependency.DependencyVersionPair{
			Name:      name,
			Version:   version,
			Ecosystem: ecosystem,
		})
	}
	assert.ElementsMatch(t, pairs, found)

	morePairs := []dependency.DependencyVersionPair{
		{
			Name:      "express",
			Version:   "4.18.2", // same as before
			Ecosystem: "npm",
		},
		{
			Name:      "requests",
			Version:   "2.31.0", // new
			Ecosystem: "pypi",
		},
	}
	results2, err := repo.ReplaceRepositoryDependencyVersions(TestCtx, repositoryID, morePairs)
	assert.NoError(t, err)
	assert.Len(t, results2, 2)

	// Confirm old association was replaced
	var count int
	err = TestDB.QueryRow(TestCtx, `SELECT COUNT(*) FROM repository_dependency_version WHERE repository_id = $1`, repositoryID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
