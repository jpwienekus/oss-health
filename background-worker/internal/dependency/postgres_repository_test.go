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
	_, err := pool.Exec(TestCtx, "TRUNCATE TABLE dependencies RESTART IDENTITY CASCADE")

	if err != nil {
		log.Fatalf("failed to truncate dependencies: %v", err)
	}

	_, err = pool.Exec(TestCtx, "TRUNCATE TABLE dependency_repository RESTART IDENTITY CASCADE")

	if err != nil {
		log.Fatalf("failed to truncate dependency_repository: %v", err)
	}
}

func SeedDependencies(pool *pgxpool.Pool) {
	_, err := pool.Exec(TestCtx, `
		INSERT INTO dependencies (name, ecosystem, github_url_resolved, github_url_resolve_failed)
		VALUES
		('react', 'npm', false, false),
		('express', 'npm', false, false),
		('flask', 'pypi', false, false)
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

func TestGetPendingDependencies(t *testing.T) {
	ClearTables(TestDB)
	SeedDependencies(TestDB)

	repository := dependency.NewPostgresRepository(TestDB)
	dependencies, err := repository.GetPendingDependencies(TestCtx, 10, 0, "npm")
	assert.NoError(t, err)
	assert.Len(t, dependencies, 2)

	names := []string{dependencies[0].Name, dependencies[1].Name}
	assert.Contains(t, names, "react")
	assert.Contains(t, names, "express")
}

func TestUpsertGithubURLs(t *testing.T) {
	ClearTables(TestDB)

	urls := []string{
		"https://github.com/facebook/react",
		"https://github.com/expressjs/express",
	}

	repository := dependency.NewPostgresRepository(TestDB)
	urlToID, err := repository.UpsertGithubURLs(TestCtx, urls)
	assert.NoError(t, err)
	assert.Len(t, urlToID, 2)

	for _, url := range urls {
		id, ok := urlToID[url]
		assert.True(t, ok, fmt.Sprintf("url %s not found in result", url))
		assert.Greater(t, id, int64(0))
	}

	// Insert duplicates again, expect same IDs returned (no duplicates)
	urlToID2, err := repository.UpsertGithubURLs(TestCtx, urls)
	assert.NoError(t, err)
	assert.Equal(t, urlToID, urlToID2)
}

func TestBatchUpdateDependencies(t *testing.T) {
	ClearTables(TestDB)
	SeedDependencies(TestDB)

	urls := []string{"https://github.com/facebook/react"}
	repository := dependency.NewPostgresRepository(TestDB)
	urlToID, err := repository.UpsertGithubURLs(TestCtx, urls)
	assert.NoError(t, err)

	deps, err := repository.GetPendingDependencies(TestCtx, 10, 0, "npm")
	assert.NoError(t, err)

	resolvedURLs := map[int64]string{}

	for _, d := range deps {
		if d.Name == "react" {
			resolvedURLs[d.ID] = "https://github.com/facebook/react"
		}
	}

	err = repository.BatchUpdateDependencies(TestCtx, deps, urlToID, resolvedURLs)
	assert.NoError(t, err)

	var resolved bool
	err = TestDB.QueryRow(TestCtx, `SELECT github_url_resolved FROM dependencies WHERE name='react'`).Scan(&resolved)

	assert.NoError(t, err)
	assert.True(t, resolved)
}

func TestMarkDependenciesAsFailed(t *testing.T) {
	ClearTables(TestDB)
	SeedDependencies(TestDB)

	repository := dependency.NewPostgresRepository(TestDB)
	deps, err := repository.GetPendingDependencies(TestCtx, 10, 0, "npm")
	assert.NoError(t, err)

	failureReasons := map[int64]string{}

	for _, d := range deps {
		failureReasons[d.ID] = "Failed to resolve URL"
	}

	err = repository.MarkDependenciesAsFailed(TestCtx, failureReasons)
	assert.NoError(t, err)

	var failed bool
	var reason string
	err = TestDB.QueryRow(TestCtx, `SELECT github_url_resolve_failed, github_url_resolve_failed_reason FROM dependencies WHERE id=$1`, deps[0].ID).Scan(&failed, &reason)

	assert.NoError(t, err)
	assert.True(t, failed)
	assert.Equal(t, "Failed to resolve URL", reason)
}
