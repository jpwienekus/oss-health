package repository

import (
	"context"
	"log"
	"sync"

	"golang.org/x/time/rate"
)

type WorkerPool struct {
	MaxParallel int
	RateLimiter *rate.Limiter
	waitGroup   sync.WaitGroup
	jobs        chan Repository
	s           *RepositoryService
}

func NewWorkerPool(maxParallel int, requestsPerSecond float64, s *RepositoryService) *WorkerPool {
	return &WorkerPool{
		MaxParallel: maxParallel,
		RateLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
		jobs:        make(chan Repository),
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

		pairs, err := workerPool.s.CloneAndParse(ctx, r)

		if err != nil {
			log.Printf("Could not parse dependencies for %s: %v", r.URL, err)
		}

		workerPool.s.ReplaceRepositoryDependencyVersions(ctx, r.ID, pairs)

		workerPool.waitGroup.Done()

	}
}

func (workerPool *WorkerPool) Submit(repository Repository) {
	workerPool.waitGroup.Add(1)
	workerPool.jobs <- repository
}

func (workerPool *WorkerPool) Wait() {
	workerPool.waitGroup.Wait()
	close(workerPool.jobs)
}
