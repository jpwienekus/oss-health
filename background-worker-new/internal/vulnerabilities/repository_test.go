package vulnerabilities_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/vulnerabilities"
)

var TestDB *pgxpool.Pool
var TestCtx = context.Background()

func ClearTables(pool *pgxpool.Pool) {
	tables := []string{
		"version_vulnerability",
		"vulnerabilities",
		"versions",
		"dependencies",
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

func TestInsertVulnerabilities(t *testing.T) {
	repo := vulnerabilities.NewVulnerabilityRepository(TestDB)
	ctx := context.Background()
	ClearTables(TestDB)

	// Seed input
	input := map[int][]vulnerabilities.Vulnerability{
		1: {
			{ID: "OSV-123", Summary: "Remote Code Execution", Severity: "HIGH"},
		},
		2: {
			{ID: "OSV-456", Summary: "XSS vulnerability", Severity: "MEDIUM"},
		},
	}

	// Insert
	inserted, err := repo.InsertVulnerabilities(ctx, input)
	assert.NoError(t, err)
	assert.Len(t, inserted, 2)

	// Query DB directly
	rows, err := TestDB.Query(ctx, `SELECT osv_id, summary, severity FROM vulnerabilities`)
	assert.NoError(t, err)

	defer rows.Close()
	found := map[string]string{}
	for rows.Next() {
		var osvID, summary, severity string
		err := rows.Scan(&osvID, &summary, &severity)
		assert.NoError(t, err)
		found[osvID] = severity
	}

	assert.Equal(t, "HIGH", found["OSV-123"])
	assert.Equal(t, "MEDIUM", found["OSV-456"])
}

func TestInsertVersionVulnerability(t *testing.T) {
	repo := vulnerabilities.NewVulnerabilityRepository(TestDB)
	ctx := context.Background()
	ClearTables(TestDB)

	// Seed dependencies and versions
	_, err := TestDB.Exec(ctx, `
		INSERT INTO dependencies (id, name, ecosystem) VALUES (1, 'lodash', 'npm');
		INSERT INTO versions (id, dependency_id, version) VALUES (10, 1, '4.17.21');
		INSERT INTO vulnerabilities (id, osv_id, summary, severity)
		VALUES (100, 'OSV-1', 'Some vuln', 'HIGH'), (200, 'OSV-2', 'Another vuln', 'LOW');
	`)
	assert.NoError(t, err)

	// Insert version-vulnerability mapping
	err = repo.InsertVersionVulnerability(ctx, map[int]int{
		100: 10,
		200: 10,
	})
	assert.NoError(t, err)

	// Validate from DB
	rows, err := TestDB.Query(ctx, `SELECT vulnerability_id, version_id FROM version_vulnerability`)
	assert.NoError(t, err)

	defer rows.Close()
	found := map[int]int{}
	for rows.Next() {
		var vulnID, versionID int
		err := rows.Scan(&vulnID, &versionID)
		assert.NoError(t, err)
		found[vulnID] = versionID
	}

	assert.Equal(t, 10, found[100])
	assert.Equal(t, 10, found[200])
}


func TestInsertVulnerabilities_Deduplication(t *testing.T) {
	repo := vulnerabilities.NewVulnerabilityRepository(TestDB)
	ctx := context.Background()
	ClearTables(TestDB)

	input := map[int][]vulnerabilities.Vulnerability{
		1: {
			{ID: "OSV-789", Summary: "Memory Leak", Severity: "MEDIUM"},
		},
	}

	// First insert
	inserted1, err := repo.InsertVulnerabilities(ctx, input)
	assert.NoError(t, err)
	assert.Len(t, inserted1, 1)

	// Second insert (duplicate)
	inserted2, err := repo.InsertVulnerabilities(ctx, input)
	assert.NoError(t, err)
	assert.Len(t, inserted2, 0) // Should not re-insert

	// Query the DB to confirm only one exists
	var count int
	err = TestDB.QueryRow(ctx, `SELECT COUNT(*) FROM vulnerabilities WHERE osv_id = 'OSV-789'`).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestInsertVersionVulnerability_InvalidVersion(t *testing.T) {
	repo := vulnerabilities.NewVulnerabilityRepository(TestDB)
	ctx := context.Background()
	ClearTables(TestDB)

	// Seed only vulnerability, not version
	_, err := TestDB.Exec(ctx, `
		INSERT INTO vulnerabilities (id, osv_id, summary, severity)
		VALUES (300, 'OSV-999', 'Invalid test', 'CRITICAL');
	`)
	assert.NoError(t, err)

	// Attempt to insert with non-existent version_id (e.g., 999)
	err = repo.InsertVersionVulnerability(ctx, map[int]int{
		300: 999, // invalid version_id
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert version vulnerabilities")
}
