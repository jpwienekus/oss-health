package repo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies"
	"github.com/oss-health/background-worker/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
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


func TestMarkScanned(t *testing.T) {
	repo := repo.NewRepoRepository(TestDB)
	ctx := context.Background()

	// Setup
	_, err := TestDB.Exec(ctx, `DELETE FROM repositories`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `DELETE FROM "user"`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `
		INSERT INTO "user" (id, github_id, github_username, access_token)
		VALUES (1, 123, 'testuser', 'token')
	`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `
		INSERT INTO repositories (id, url, github_id, user_id, scan_day, scan_hour, scan_status)
		VALUES (1, 'http://example.com', 123, 1, 1, 1, 'pending')
	`)
	assert.NoError(t, err)

	// Act
	err = repo.MarkScanned(ctx, 1)
	assert.NoError(t, err)

	// Assert
	var status string
	err = TestDB.QueryRow(ctx, `SELECT scan_status FROM repositories WHERE id = 1`).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "done", status)
}

func TestMarkFailed(t *testing.T) {
	repo := repo.NewRepoRepository(TestDB)
	ctx := context.Background()

	// Setup
	_, err := TestDB.Exec(ctx, `DELETE FROM repositories`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `DELETE FROM "user"`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `
		INSERT INTO "user" (id, github_id, github_username, access_token)
		VALUES (1, 456, 'failuser', 'token')
	`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `
		INSERT INTO repositories (id, url, github_id, user_id, scan_status)
		VALUES (2, 'http://fail.com', 456, 1, 'pending')
	`)
	assert.NoError(t, err)

	// Act
	failMsg := "fail message"
	err = repo.MarkFailed(ctx, 2, failMsg)
	assert.NoError(t, err)

	// Assert
	var status, dbFailMsg string
	err = TestDB.QueryRow(ctx, `SELECT scan_status, error_message FROM repositories WHERE id = 2`).Scan(&status, &dbFailMsg)
	assert.NoError(t, err)
	assert.Equal(t, "error", status)
	assert.Equal(t, failMsg, dbFailMsg)
}


func TestInsertAndDeleteRepositoryDependencyVersions(t *testing.T) {
	repo := repo.NewRepoRepository(TestDB)
	ctx := context.Background()

	ClearTables(TestDB)

	// Setup required records
	_, err := TestDB.Exec(ctx, `INSERT INTO "user" (id, github_id, github_username, access_token) VALUES (1, 999, 'depuser', 'token')`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `INSERT INTO repositories (id, url, github_id, user_id, scan_status) VALUES (1, 'https://example.com/repo', 999, 1, 'pending')`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `INSERT INTO dependencies (id, name, ecosystem) VALUES (1, 'dep1', 'npm'), (2, 'dep2', 'npm')`)
	assert.NoError(t, err)

	_, err = TestDB.Exec(ctx, `INSERT INTO versions (id, dependency_id, version) VALUES (10, 1, '1.0.0'), (20, 2, '2.0.0')`)
	assert.NoError(t, err)

	// Act: Insert
	deps := []dependencies.DependencyVersion{
		{ID: 10, DependencyId: 1},
		{ID: 20, DependencyId: 2},
	}
	err = repo.InsertRepositoryDependencyVersion(ctx, 1, deps)
	assert.NoError(t, err)

	// Assert insert
	var count int
	err = TestDB.QueryRow(ctx, `SELECT COUNT(*) FROM repository_dependency_version WHERE repository_id = 1`).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	// Act: Delete
	err = repo.DeleteRepositoryDependencyVersion(ctx, 1)
	assert.NoError(t, err)

	// Assert delete
	err = TestDB.QueryRow(ctx, `SELECT COUNT(*) FROM repository_dependency_version WHERE repository_id = 1`).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
