package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/dependency/resolvers"
	"github.com/oss-health/background-worker/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("DATABASE_URL")

	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)

	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer db.Close()
	repository := dependency.NewPostgresRepository(db)
	rateLimiter := utils.NewDefaultRateLimiter()
	dependency.InitRateLimiters(rateLimiter)
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
