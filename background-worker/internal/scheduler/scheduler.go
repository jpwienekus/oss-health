package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/oss-health/background-worker/internal/utils"
	"github.com/oss-health/background-worker/pkg/fetcher"
	"github.com/robfig/cron/v3"
	"golang.org/x/time/rate"
)

const (
	NpmRps    = 10
	NpmBurst  = 10
	PypiRps   = 1
	PypiBurst = 1
)

func Start() {
	rateLimiter := utils.NewDefaultRateLimiter()
	initRateLimiters(rateLimiter)

	c := cron.New(
		cron.WithSeconds(),
	)

	buffer := 10
	npmRequestCapability := (NpmRps * 60) - buffer
	pypiRequestCapability := (PypiRps * 60) - buffer

	_, err := c.AddFunc("0 */1 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		log.Println("Starting scheduled fetch job: npm")
		if err := fetcher.ResolvePendingDependencies(ctx, npmRequestCapability, 0, "npm", rateLimiter, fetcher.Resolvers); err != nil {
			log.Printf("Error running npm fetch job: %v", err)
		} else {
			log.Printf("Finished batch for npm")
		}
	})

	_, err = c.AddFunc("0 */1 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		log.Println("Starting scheduled fetch job: pypi")
		if err := fetcher.ResolvePendingDependencies(ctx, pypiRequestCapability, 0, "pypi", rateLimiter, fetcher.Resolvers); err != nil {
			log.Printf("Error running pypi fetch job: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("failed to schedule tasks: %v", err)
	}

	log.Println("Scheduler started")
	c.Start()
}

func initRateLimiters(rateLimiter *utils.DefaultRateLimiter) {
	registerRateLimter(rateLimiter, "npm", NpmRps, NpmBurst)
	registerRateLimter(rateLimiter, "pypi", PypiRps, PypiBurst)
}

func registerRateLimter(rateLimiter *utils.DefaultRateLimiter, registry string, rps int, burst int) {
	periodPerRequest := time.Second / time.Duration(rps)
	rateLimiter.RegisterLimiter(registry, rate.Every(periodPerRequest), burst)
}
