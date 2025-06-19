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
	// fetcher.ResolvePendingDependencies(ctx, 120, 0, "npm")
	scheduler.Start()
	log.Println("Done witih batch")

	// Graceful shutdown block
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down gracefully.")
}

func initRateLimiters() {
	// Allow 1 request per second (60/min) with small burst
	utils.RegisterLimiter("npm", rate.Every(time.Second), 2)

	// Allow 1 request every 2 seconds (30/min)
	utils.RegisterLimiter("pypi", rate.Every(2*time.Second), 1)
}
