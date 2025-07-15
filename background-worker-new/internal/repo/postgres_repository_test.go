package repo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/oss-health/background-worker/internal/db"
	// "github.com/oss-health/background-worker/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/stretchr/testify/assert"
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
