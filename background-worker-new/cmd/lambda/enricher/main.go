package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies"
	"github.com/oss-health/background-worker/internal/vulnerabilities"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	dependencyService := vulnerabilities.NewVulnerabilityService(
		vulnerabilities.NewVulnerabilityRepository(db),
		dependencies.NewDependencyRepository(db),
		vulnerabilities.NewOSVClient(),
	)

	if err := dependencyService.RunOSVScan(ctx); err != nil {
		log.Fatalf("Failed to run OSV scan: %v", err)
	}
}
