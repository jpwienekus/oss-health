package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oss-health/background-worker/db"
	// "github.com/oss-health/background-worker/fetcher"
	"github.com/oss-health/background-worker/scheduler"
	"github.com/oss-health/background-worker/utils"
	"golang.org/x/time/rate"
)

func main() {
	ctx := context.Background()
	connStr := "postgres://dev-user:password@localhost:5432/dev_db"

	if err := db.Connect(ctx, connStr); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	initRateLimiters()
	// fetcher.ResolvePendingDependencies(ctx, 10, 0, "npm")
	// log.Println("Done witih batch")
	scheduler.Start()

	// Graceful shutdown block
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down gracefully.")
}

func initRateLimiters() {
	npmBurst := 10
	npmRequestsPerSecond := 10
	periodPerRequest := time.Second / time.Duration(npmRequestsPerSecond)
	utils.RegisterLimiter("npm", rate.Every(periodPerRequest), npmBurst)

	// Allow 1 request every 2 seconds (30/min)
	utils.RegisterLimiter("pypi", rate.Every(2*time.Second), 1)
}
