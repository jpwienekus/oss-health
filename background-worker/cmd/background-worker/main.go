package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/scheduler"
)

func main() {
	ctx := context.Background()
	connStr := "postgres://dev-user:password@localhost:5432/dev_db"

	if err := db.Connect(ctx, connStr); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	scheduler.Start()

	// Graceful shutdown block
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down gracefully.")
}

