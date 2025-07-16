package repo

import (
	"context"
	"log"
	"sync"

	"github.com/oss-health/background-worker/internal/db"
	"golang.org/x/time/rate"
)

type WorkerPool struct {
	MaxParallel int
	RateLimiter *rate.Limiter
	waitGroup   sync.WaitGroup
	jobs        chan db.Repository
	s           *RepositoryService
}

func NewWorkerPool(maxParallel int, requestsPerSecond float64, s *RepositoryService) *WorkerPool {
	return &WorkerPool{
		MaxParallel: maxParallel,
		RateLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
		jobs:        make(chan db.Repository),
		s:           s,
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

		dependencies, err := workerPool.s.ExtractDependencies(ctx, r)
		if err != nil {
			log.Printf("[worker %d] failed to extract dependencies for %s: %v", id, r.URL, err)
		}

		if err := workerPool.s.ReplaceRepositoryDependencyVersions(ctx, r.ID, dependencies, id); err != nil {
			log.Printf("[worker %d] failed to replace dependency versions for repository ID %d: %v", id, r.ID, err)
		}

		workerPool.waitGroup.Done()

	}
}

func (workerPool *WorkerPool) Submit(repository db.Repository) {
	workerPool.waitGroup.Add(1)
	workerPool.jobs <- repository
}

func (workerPool *WorkerPool) Wait() {
	workerPool.waitGroup.Wait()
	close(workerPool.jobs)
}
