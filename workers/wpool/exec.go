package wpool

import (
	"context"
	"fmt"
	"sync"
)

func Worker[A any, T any](ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job[A, T], results chan<- Result[T]) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			// fan-in job execution multiplexing results into the results channel
			results <- job.Execute(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result[T]{
				Err: ctx.Err(),
			}
			return
		}
	}
}

type WorkerPool[A any, T any] struct {
	workersCount int
	jobs         chan Job[A, T]
	results      chan Result[T]
	Done         chan struct{}
}

func New[A any, T any](wcount int) WorkerPool[A, T] {
	return WorkerPool[A, T]{
		workersCount: wcount,
		jobs:         make(chan Job[A, T], wcount),
		results:      make(chan Result[T], wcount),
		Done:         make(chan struct{}),
	}
}

func (wp WorkerPool[A, T]) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		// fan out worker goroutines
		//reading from jobs channel and
		//pushing calcs into results channel
		go Worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp WorkerPool[A, T]) Results() <-chan Result[T] {
	return wp.results
}

func (wp WorkerPool[A, T]) GenerateFrom(jobsBulk []Job[A, T]) {
	for i := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}
