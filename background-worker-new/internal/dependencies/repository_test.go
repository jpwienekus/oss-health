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
