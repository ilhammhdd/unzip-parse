package entities

import (
	"sync"
)

type WorkerPool struct {
	Job          chan func() interface{}
	Result       chan interface{}
	WorkerAmount int
	wg           sync.WaitGroup
}

func NewPool(workerAmount int) *WorkerPool {
	p := &WorkerPool{
		Job:          make(chan func() interface{}),
		Result:       make(chan interface{}),
		WorkerAmount: workerAmount,
	}

	return p
}

func (p *WorkerPool) Work(job func() interface{}) {
	p.wg.Add(p.WorkerAmount)

	for i := 0; i < p.WorkerAmount; i++ {
		GoSafely(func() {
			p.runWorker()
			p.wg.Done()
		})
	}

	for j := 0; j < p.WorkerAmount; j++ {
		p.Job <- job
	}
	close(p.Job)

	GoSafely(func() {
		p.wg.Wait()
		close(p.Result)
	})
}

func (p *WorkerPool) runWorker() {
	for j := range p.Job {
		p.Result <- j()
	}
}
