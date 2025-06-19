package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/oss-health/background-worker/fetcher"
)

func Start() {
	c := cron.New(
		cron.WithSeconds(),
	)

	// Schedule: Every 1 minute
	_, err := c.AddFunc("0 */1 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		log.Println("Starting scheduled fetch job: npm")
		if err := fetcher.ResolvePendingDependencies(ctx, 50, 0, "npm"); err != nil {
			log.Printf("Error running npm fetch job: %v", err)
		} else {
			log.Printf("Finished batch for npm")
		}
	})

	// // Schedule: Every 2 minutes
	// _, err = c.AddFunc("0 */2 * * * *", func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	// 	defer cancel()
	//
	// 	log.Println("⏳ Starting scheduled fetch job: pypi")
	// 	if err := fetcher.ResolvePendingDependencies(ctx, 50, 0, "pypi"); err != nil {
	// 		log.Printf("Error running pypi fetch job: %v", err)
	// 	}
	// })

	if err != nil {
		log.Fatalf("failed to schedule tasks: %v", err)
	}

	log.Println("Scheduler started")
	c.Start()
}
