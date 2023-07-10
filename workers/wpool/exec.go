package wpool

import (
	"context"
	"fmt"
	"sync"
)

// Why use a wait group, when we can just have a done channel?

func Worker[A any, T any](ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job[A, T], results chan<- Result[T]) {
	defer wg.Done()
	fmt.Printf("🤖 Worker starting ...\n")
	for {
		// THIS IS DEADLOCK BLOCKING IF THERE IS NO WORK. WE NEED TO PRE-EMPT WITH A TIMEOUT
		//
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("🧊 Worker DONE(remote closed) ...\n")
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

	// The current design is not very good on startup. No consideration is
	// taken of the fact that there maybe no work to do.
	//
	fmt.Printf("🦄 Running %v workers\n", wp.workersCount)

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		// fan out worker goroutines
		// reading from jobs channel and
		// pushing calcs into results channel

		go Worker(ctx, &wg, wp.jobs, wp.results)
	}

	// sometimes, not all results are read from the results channel before the
	// program finishes.
	//
	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp WorkerPool[A, T]) Results() <-chan Result[T] {
	return wp.results
}

func (wp WorkerPool[A, T]) Dispatch(jobsBulk []Job[A, T]) {
	for i := range jobsBulk {
		wp.jobs <- jobsBulk[i] // THIS BLOCKS WHEN THE CHANNEL IS FULL
	}

}

func (wp WorkerPool[A, T]) Finish() {
	fmt.Printf("💤 Closing job stream.\n")

	// shouldn't close the jobs channel here, it the client responsibility
	// to do so.
	//
	// THIS IS A PROBLEM. WE CANT HAVE A STREAMING MODEL WITH THIS CLOSURE
	//
	// We need to decouple the close from dispatching of jobs to the pool
	//
	close(wp.jobs)
}
