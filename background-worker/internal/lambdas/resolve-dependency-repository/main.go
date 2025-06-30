package main

import (
	"context"
	"log"
	"time"
	"sync"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/dependency/resolvers"
	"github.com/oss-health/background-worker/internal/utils"
)

func main() {
	// connectionString := "postgres://dev-user:password@localhost:5432/dev_db"
	connectionString := "postgresql://postgres.gfpivacysduostopkekw:4-dzBCK8Ptyg.FTukiBB@aws-0-eu-central-1.pooler.supabase.com:5432/postgres"

	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)

	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer db.Close()
	repository := dependency.NewPostgresRepository(db)
	rateLimiter := utils.NewDefaultRateLimiter()
	service := dependency.NewDependencyService(repository, rateLimiter, resolvers.Resolvers)

	buffer := 10
	npmRequestCapability := (dependency.NpmRps * 60) - buffer
	// pypiRequestCapability := (dependency.PypiRps * 60) - buffer

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		resolvePendingDependencies("npm", npmRequestCapability, service)
	}()

	// go func() {
	// 	defer wg.Done()
	// 	resolvePendingDependencies("pypi", pypiRequestCapability, service)
	// }()

	wg.Wait()
}

func resolvePendingDependencies(ecosystem string, batchSize int, service *dependency.DependencyService) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	err := service.ResolvePendingDependencies(ctx, batchSize, 0, ecosystem)

	if err != nil {
		log.Printf("Error running %s fetch job: %v", ecosystem, err)
	}
}
