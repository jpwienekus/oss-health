package scanner

import (
	"context"
	"log"
	"sync"

	"github.com/oss-health/background-worker/internal/repository"
	"golang.org/x/time/rate"
)

type WorkerPool struct {
	MaxParallel int
	RateLimiter *rate.Limiter
	waitGroup   sync.WaitGroup
	jobs        chan repository.Repository
}

func NewWorkerPool(maxParallel int, requestsPerSecond float64) *WorkerPool {
	return &WorkerPool{
		MaxParallel: maxParallel,
		RateLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
		jobs:        make(chan repository.Repository),
	}
}

func (workerPool *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < workerPool.MaxParallel; i++ {
		go workerPool.worker(i, ctx)

	}
}

func (workerPool *WorkerPool) worker(id int, ctx context.Context) {
	for r := range workerPool.jobs {
		err := workerPool.RateLimiter.Wait(context.Background())

		if err != nil {
			log.Printf("[worker %d] rate limiter error: %v", id, err)
			continue
		}

		log.Printf("[worker %d] processing %s", id, r.URL)
		CloneAndParse(ctx, r)

		workerPool.waitGroup.Done()

	}
}

func (workerPool *WorkerPool) Submit(repository repository.Repository) {
	workerPool.waitGroup.Add(1)
	workerPool.jobs <- repository
}

func (workerPool *WorkerPool) Wait() {
	workerPool.waitGroup.Wait()
	close(workerPool.jobs)
}
