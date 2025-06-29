package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/utils"
	"github.com/oss-health/background-worker/internal/dependency/resolvers"

	"github.com/robfig/cron/v3"
)

func Start() {
	// TODO: env var
	connectionString := "postgres://dev-user:password@localhost:5432/dev_db"

	ctx := context.Background()
	db, err := db.Connect(ctx, connectionString)

	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer db.Close()

	rateLimiter := utils.NewDefaultRateLimiter()
	dependency.InitRateLimiters(rateLimiter)

	repository := dependency.NewPostgresRepository(db)
	service := dependency.NewDependencyService(repository, rateLimiter, resolvers.Resolvers)

	buffer := 10
	npmRequestCapability := (dependency.NpmRps * 60) - buffer
	pypiRequestCapability := (dependency.PypiRps * 60) - buffer

	StartScheduler(service, npmRequestCapability, pypiRequestCapability)
}

func StartScheduler(service *dependency.DependencyService, npmRequestCapability int, pypiRequestCapability int) {
	rateLimiter := utils.NewDefaultRateLimiter()
	dependency.InitRateLimiters(rateLimiter)

	c := cron.New(cron.WithSeconds())

	scheduleFetchJob(c, "npm", npmRequestCapability, service)
	scheduleFetchJob(c, "pypi", pypiRequestCapability, service)

	log.Println("Scheduler started")
	c.Start()
}

func scheduleFetchJob(c *cron.Cron, ecosystem string, batchSize int, service *dependency.DependencyService) {
	spec := "0 */1 * * * *"
	_, err := c.AddFunc(spec, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		log.Printf("Starting scheduled fetch job: %s", ecosystem)
		if err := service.ResolvePendingDependencies(ctx, batchSize, 0, ecosystem); err != nil {
			log.Printf("Error running %s fetch job: %v", ecosystem, err)
		} else {
			log.Printf("Finished batch for %s", ecosystem)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule %s job: %v", ecosystem, err)
	}
}
