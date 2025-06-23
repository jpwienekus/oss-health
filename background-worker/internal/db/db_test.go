package db_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/stretchr/testify/assert"
)

var testCtx = context.Background()

func TestMain(m *testing.M) {
	connStr := "postgres://test-user:password@localhost:5434/test_db"
	err := db.Connect(testCtx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}

	clearTables()

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func clearTables() {
	_, err := db.Pool.Exec(testCtx, "TRUNCATE TABLE dependencies RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("failed to truncate dependencies: %v", err)
	}

	_, err = db.Pool.Exec(testCtx, "TRUNCATE TABLE dependency_repository RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("failed to truncate dependency_repository: %v", err)
	}
}

func seedDependencies() {
	_, err := db.Pool.Exec(testCtx, `
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

func TestGetPendingDependencies(t *testing.T) {
	clearTables()
	seedDependencies()

	dependencies, err := dependency.GetPendingDependencies(testCtx, 10, 0, "npm")
	assert.NoError(t, err)
	assert.Len(t, dependencies, 2)

	names := []string{dependencies[0].Name, dependencies[1].Name}
	assert.Contains(t, names, "react")
	assert.Contains(t, names, "express")
}


func TestUpsertGithubURLs(t *testing.T) {
	clearTables()

	urls := []string{
		"https://github.com/facebook/react",
		"https://github.com/expressjs/express",
	}

	urlToID, err := dependency.UpsertGithubURLs(testCtx, urls)
	assert.NoError(t, err)
	assert.Len(t, urlToID, 2)

	for _, url := range urls {
		id, ok := urlToID[url]
		assert.True(t, ok, fmt.Sprintf("url %s not found in result", url))
		assert.Greater(t, id, int64(0))
	}

	// Insert duplicates again, expect same IDs returned (no duplicates)
	urlToID2, err := dependency.UpsertGithubURLs(testCtx, urls)
	assert.NoError(t, err)
	assert.Equal(t, urlToID, urlToID2)
}

func TestBatchUpdateDependencies(t *testing.T) {
	clearTables()
	seedDependencies()

	urls := []string{"https://github.com/facebook/react"}
	urlToID, err := dependency.UpsertGithubURLs(testCtx, urls)
	assert.NoError(t, err)

	deps, err := dependency.GetPendingDependencies(testCtx, 10, 0, "npm")
	assert.NoError(t, err)

	resolvedURLs := map[int64]string{}

	for _, d := range deps {
		if d.Name == "react" {
			resolvedURLs[d.ID] = "https://github.com/facebook/react"
		}
	}

	err = dependency.BatchUpdateDependencies(testCtx, deps, urlToID, resolvedURLs)
	assert.NoError(t, err)

	var resolved bool
	err = db.Pool.QueryRow(testCtx, `SELECT github_url_resolved FROM dependencies WHERE name='react'`).Scan(&resolved)

	assert.NoError(t, err)
	assert.True(t, resolved)
}

func TestMarkDependenciesAsFailed(t *testing.T) {
	clearTables()
	seedDependencies()

	deps, err := dependency.GetPendingDependencies(testCtx, 10, 0, "npm")
	assert.NoError(t, err)

	failureReasons := map[int64]string{}

	for _, d := range deps {
		failureReasons[d.ID] = "Failed to resolve URL"
	}

	err = dependency.MarkDependenciesAsFailed(testCtx, failureReasons)
	assert.NoError(t, err)

	var failed bool
	var reason string
	err = db.Pool.QueryRow(testCtx, `SELECT github_url_resolve_failed, github_url_resolve_failed_reason FROM dependencies WHERE id=$1`, deps[0].ID).Scan(&failed, &reason)

	assert.NoError(t, err)
	assert.True(t, failed)
	assert.Equal(t, "Failed to resolve URL", reason)
}
