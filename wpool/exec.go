package wpool

import (
	"context"
	"fmt"
)

func worker(ctx context.Context, jobs <-chan Job, results chan<- Result, sem semaphore) {
	defer sem.release()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			// fan-in job execution multiplexing results into the results channel
			results <- job.execute(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}

type semaphore interface {
	acquire()
	release()
	wait()
	close()
}

type slot struct{}
type slots chan slot

type execution struct {
	slots slots
}

func newExecutionSlots(capacity int) execution {
	slots := make(chan slot, capacity)
	return execution{slots: slots}
}

func (e execution) acquire() {
	e.slots <- slot{}
}

func (e execution) release() {
	<-e.slots
}

func (e execution) wait() {
	for i := 0; i < cap(e.slots); i++ {
		e.slots <- slot{}
	}
}

func (e execution) close() {
	close(e.slots)
}

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
}

func New(wcount int) WorkerPool {
	return WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		results:      make(chan Result, wcount),
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	eSlots := newExecutionSlots(wp.workersCount)
	defer eSlots.close()

	for i := 0; i < wp.workersCount; i++ {
		eSlots.acquire()
		// fan out worker goroutines
		//reading from jobs channel and
		//pushing calcs into results channel
		go worker(ctx, wp.jobs, wp.results, eSlots)
	}

	eSlots.wait()
	close(wp.results)
}

func (wp WorkerPool) Results() <-chan Result {
	return wp.results
}

func (wp WorkerPool) GenerateFrom(jobsBulk []Job) {
	for i := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}
