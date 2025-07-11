package storage_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/parser"
	"github.com/oss-health/background-worker/storage"
)

func TestGetRepositoriesBySchedule(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect the query and provide rows
	mock.ExpectQuery(`SELECT id, url FROM repositories WHERE scan_day = \$1 AND scan_hour = \$2 ORDER BY id`).
		WithArgs(2, 14).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url"}).
			AddRow(1, "https://github.com/oss/repo1").
			AddRow(2, "https://github.com/oss/repo2"))

	store := &storage.PostgresStore{DB: db}

	repos, err := store.GetRepositoriesBySchedule(2, 14)
	assert.NoError(t, err)
	assert.Len(t, repos, 2)
	assert.Equal(t, 1, repos[0].ID)
	assert.Equal(t, "https://github.com/oss/repo1", repos[0].URL)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkRepositoryAsScanned(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`(?s)UPDATE repositories\s+SET\s+last_scanned_at = now\(\),\s+scan_status = 'done'\s+WHERE id = \$1`).
		WithArgs(123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := &storage.PostgresStore{DB: db}
	err = store.MarkRepositoryAsScanned(123)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkRepositoryAsFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`(?s)UPDATE repositories\s+SET\s+scan_status = 'error',\s+error_message = \$1\s+WHERE id = \$2`).
		WithArgs("failed to parse repo", 456).
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := &storage.PostgresStore{DB: db}
	err = store.MarkRepositoryAsFailed(456, "failed to parse repo")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkRepositoryAsScanned_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`(?s)UPDATE repositories.*scan_status = 'done'.*`).
		WithArgs(999).
		WillReturnError(fmt.Errorf("db is down"))

	store := &storage.PostgresStore{DB: db}
	err = store.MarkRepositoryAsScanned(999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to mark repository as scanned")
}

func TestDeleteRepositoryDependencyVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`(?s)DELETE FROM repository_dependency_version\s+WHERE repository_id = \$1`).
		WithArgs(321).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 0 insert ID, 1 row affected

	store := &storage.PostgresStore{DB: db}
	err = store.DeleteRepositoryDependencyVersion(321)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetExistingDependencies(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	deps := []parser.Dependency{
		{Name: "lodash", Ecosystem: "npm"},
		{Name: "requests", Ecosystem: "pypi"},
	}

	mock.ExpectQuery(`(?s)SELECT id, name, ecosystem FROM dependencies WHERE .*`).
		WithArgs("lodash", "npm", "requests", "pypi").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "ecosystem"}).
			AddRow(1, "lodash", "npm").
			AddRow(2, "requests", "pypi"))

	store := &storage.PostgresStore{DB: db}
	result, err := store.GetExistingDependencies(deps)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "lodash", result[0].Name)
	assert.Equal(t, "pypi", result[1].Ecosystem)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertDependencies(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	deps := []parser.Dependency{
		{Name: "express", Ecosystem: "npm"},
		{Name: "flask", Ecosystem: "pypi"},
	}

	// Match INSERT INTO ... VALUES ($1, $2), ($3, $4)
	mock.ExpectExec(`(?s)INSERT INTO dependencies \(name, ecosystem\)\s+VALUES\s+\(\$1, \$2\),\s+\(\$3, \$4\)`).
		WithArgs("express", "npm", "flask", "pypi").
		WillReturnResult(sqlmock.NewResult(0, 2)) // 0 insert ID, 2 rows affected

	store := &storage.PostgresStore{DB: db}
	err = store.InsertDependencies(deps)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}


