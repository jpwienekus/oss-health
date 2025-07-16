package dependencies_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies"
	"github.com/oss-health/background-worker/internal/repo/parsers"
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

func TestInsertAndGetDependencies(t *testing.T) {
	repo := dependencies.NewDependencyRepository(TestDB)
	ctx := context.Background()

	ClearTables(TestDB)

	// Insert
	toInsert := []parsers.Dependency{
		{Name: "lodash", Ecosystem: "npm"},
		{Name: "express", Ecosystem: "npm"},
	}
	err := repo.InsertDependencies(ctx, toInsert)
	assert.NoError(t, err)

	// Get
	existing, err := repo.GetExistingDependencies(ctx, toInsert)
	assert.NoError(t, err)
	assert.Len(t, existing, 2)

	names := map[string]bool{}
	for _, d := range existing {
		names[d.Name] = true
	}
	assert.True(t, names["lodash"])
	assert.True(t, names["express"])
}

func TestInsertAndGetDependencyVersions(t *testing.T) {
	repo := dependencies.NewDependencyRepository(TestDB)
	ctx := context.Background()

	ClearTables(TestDB)

	// Seed base dependencies
	_, err := TestDB.Exec(ctx, `INSERT INTO dependencies (id, name, ecosystem) VALUES (1, 'lodash', 'npm'), (2, 'express', 'npm')`)
	assert.NoError(t, err)

	// Insert versions
	toInsert := []db.DependencyVersion{
		{Version: "1.0.0", DependencyId: 1},
		{Version: "4.17.1", DependencyId: 2},
	}
	err = repo.InsertDependencyVersions(ctx, toInsert)
	assert.NoError(t, err)

	// Get versions
	input := []parsers.Dependency{
		{Name: "lodash", Ecosystem: "npm", Version: "1.0.0"},
		{Name: "express", Ecosystem: "npm", Version: "4.17.1"},
	}
	versions, err := repo.GetExistingVersionsForDependencies(ctx, input)
	assert.NoError(t, err)
	assert.Len(t, versions, 2)

	versionMap := map[string]string{}
	for _, v := range versions {
		versionMap[v.DependencyName] = v.Version
	}

	assert.Equal(t, "1.0.0", versionMap["lodash"])
	assert.Equal(t, "4.17.1", versionMap["express"])
}

func TestDependencyUrlResolutionWorkflow(t *testing.T) {
	repo := dependencies.NewDependencyRepository(TestDB)
	ctx := context.Background()

	ClearTables(TestDB)

	_, err := TestDB.Exec(ctx, `INSERT INTO "user" (id, github_id, github_username, access_token) VALUES (1, 999, 'depuser', 'token')`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `INSERT INTO repositories (id, url, github_id, user_id, scan_status) VALUES (1, 'https://example.com/repo', 999, 1, 'pending')`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `INSERT INTO dependencies (id, name, ecosystem) VALUES (1, 'dep1', 'npm'), (2, 'dep2', 'npm')`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `INSERT INTO versions (id, dependency_id, version) VALUES (10, 1, '1.0.0'), (20, 2, '2.0.0')`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `
		INSERT INTO repository_dependency_version (repository_id, dependency_id, version_id)
		VALUES 
		(1, 1, 10),
		(1, 2, 20)
	`)
	assert.NoError(t, err)

	pending, err := repo.GetDependenciesPendingUrlResolution(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, pending, 2)

	// Step 6: Insert into dependency_repository
	urlMap := map[int]string{
		1: "https://github.com/axios/axios",
		2: "https://github.com/facebook/react",
	}
	err = repo.InsertDependencyRepository(ctx, urlMap)
	assert.NoError(t, err)

	// Step 7: Fetch dependency_repository IDs
	urlIDs, err := repo.GetDependencyRepositoryIds(ctx, urlMap)
	assert.NoError(t, err)
	assert.Len(t, urlIDs, 2)

	// Step 8: Map dependency_id -> dependency_repository_id
	depRepoMap := make(map[int]int)
	for depID, url := range urlMap {
		repoID, ok := urlIDs[url]
		assert.True(t, ok, "expected repo ID for URL %s", url)
		depRepoMap[depID] = repoID
	}

	// Step 9: Update dependencies to point to dependency_repository
	err = repo.UpdateDependencyRepositoryIDs(ctx, depRepoMap)
	assert.NoError(t, err)

	// Step 10: Check again — no unresolved dependencies now
	pendingAfterUpdate, err := repo.GetDependenciesPendingUrlResolution(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, pendingAfterUpdate, 0)
}

func TestMarkDependenciesAsFailed(t *testing.T) {
	repo := dependencies.NewDependencyRepository(TestDB)
	ctx := context.Background()

	ClearTables(TestDB)

	_, err := TestDB.Exec(ctx, `
		INSERT INTO dependencies (id, name, ecosystem)
		VALUES (1, 'vue', 'npm'), (2, 'jquery', 'npm')
	`)
	assert.NoError(t, err)

	failureReasons := map[int]string{
		1: "Invalid GitHub URL",
		2: "Timeout fetching metadata",
	}

	err = repo.MarkDependenciesAsFailed(ctx, failureReasons)
	assert.NoError(t, err)

	rows, err := TestDB.Query(ctx, `SELECT id, error_message FROM dependencies WHERE id IN (1, 2)`)
	assert.NoError(t, err)

	reasonMap := map[int]string{}
	for rows.Next() {
		var id int
		var reason string
		err := rows.Scan(&id, &reason)
		assert.NoError(t, err)
		reasonMap[id] = reason
	}
	rows.Close()

	assert.Equal(t, "Invalid GitHub URL", reasonMap[1])
	assert.Equal(t, "Timeout fetching metadata", reasonMap[2])
}
