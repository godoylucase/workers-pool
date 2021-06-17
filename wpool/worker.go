package wpool

import "sync"

type worker struct {
	wg      *sync.WaitGroup
	jobs    <-chan job
	results chan<- result
}

func (w worker) execute() {
	for job := range w.jobs {
		w.results <- job.Execute()
	}
	w.wg.Done()
}

type Pool struct {
	workersCount int
	jobs         <-chan job
	results      chan result
	done         chan interface{}
}

func New(workersCount int, jobs <-chan job) Pool {
	results := make(chan result, cap(jobs))

	return Pool{
		workersCount: workersCount,
		jobs:         jobs,
		results:      results,
	}
}

func (p Pool) Start() <-chan result {
	var wg sync.WaitGroup
	defer close(p.results)

	for i := 0; i < p.workersCount; i++ {
		wg.Add(1)
		go worker{
			wg:      &wg,
			jobs:    p.jobs,
			results: p.results,
		}.execute()
	}
	wg.Wait()

	return p.results
}
