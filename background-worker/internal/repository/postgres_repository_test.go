package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var TestCtx = context.Background()

func TestGetRepositoriesForDay(t *testing.T) {
	connStr := "postgres://test-user:password@localhost:5434/test_db"
	pool, err := db.Connect(TestCtx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	r := repository.NewRepositoryRepository(pool)
	ctx := context.Background()

	_, err = pool.Exec(ctx, `DELETE FROM repositories`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `DELETE FROM "user"`)
	require.NoError(t, err)

	// Setup
	_, err = pool.Exec(ctx, `
		INSERT INTO "user" (id, github_id, github_username, access_token)
		VALUES 
			(1, 101, 'user101', 'token1'),
			(2, 102, 'user102', 'token2'),
			(3, 103, 'user103', 'token3')
	`)
	require.NoError(t, err)

	now := time.Now()
	_, err = pool.Exec(ctx, `
		INSERT INTO repositories (id, url, github_id, user_id, last_scanned_at, scan_status, scan_day, scan_hour)
		VALUES
			(1, 'http://example.com', 101, 1, $1, 'completed', 2, 14),
			(2, 'http://example2.com', 102, 2, $2, 'pending', 1, 3),
			(3, 'http://old.com', 103, 3, $3, 'completed', 5, 22)
	`, now, now, now)
	require.NoError(t, err)

	// Act
	repos, err := r.GetRepositoriesForDay(ctx, 2, 14)
	require.NoError(t, err)

	// Assert
	assert.Len(t, repos, 1)

	assert.Equal(t, "http://example.com", repos[0].URL)
	assert.NotNil(t, repos[0].LastScannedAt)
	assert.Equal(t, "completed", repos[0].ScanStatus)

	_, err = pool.Exec(ctx, `DELETE FROM repositories`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `DELETE FROM "user"`)
	require.NoError(t, err)
}

func TestMarkScanned(t *testing.T) {
	connStr := "postgres://test-user:password@localhost:5434/test_db"
	pool, err := db.Connect(context.Background(), connStr)
	require.NoError(t, err)
	defer pool.Close()

	repo := repository.NewRepositoryRepository(pool)
	ctx := context.Background()

	// Setup
	_, err = pool.Exec(ctx, `DELETE FROM repositories`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `DELETE FROM "user"`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO "user" (id, github_id, github_username, access_token)
		VALUES (1, 123, 'testuser', 'token')
	`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO repositories (id, url, github_id, user_id, scan_status)
		VALUES (1, 'http://example.com', 123, 1, 'pending')
	`)
	require.NoError(t, err)

	// Act
	repo.MarkScanned(ctx, 1)

	// Assert
	var status string
	err = pool.QueryRow(ctx, `SELECT scan_status FROM repositories WHERE id = 1`).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "done", status)
}

func TestMarkFailed(t *testing.T) {
	connStr := "postgres://test-user:password@localhost:5434/test_db"
	pool, err := db.Connect(context.Background(), connStr)
	require.NoError(t, err)
	defer pool.Close()

	repo := repository.NewRepositoryRepository(pool)
	ctx := context.Background()

	// Setup
	_, err = pool.Exec(ctx, `DELETE FROM repositories`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `DELETE FROM "user"`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO "user" (id, github_id, github_username, access_token)
		VALUES (1, 456, 'failuser', 'token')
	`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO repositories (id, url, github_id, user_id, scan_status)
		VALUES (2, 'http://fail.com', 456, 1, 'pending')
	`)
	require.NoError(t, err)

	// Act
	failMsg := "fail message"
	repo.MarkFailed(ctx, 2, failMsg)
	require.NoError(t, err)

	// Assert
	var status, dbFailMsg string
	err = pool.QueryRow(ctx, `SELECT scan_status, error_message FROM repositories WHERE id = 2`).Scan(&status, &dbFailMsg)
	require.NoError(t, err)
	assert.Equal(t, "error", status)
	assert.Equal(t, failMsg, dbFailMsg)
}
