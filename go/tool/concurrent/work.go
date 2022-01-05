package concurrent

import (
	"sync"
)

type Worker interface {
	Task()
}

type WorkerPool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func NewWork(maxGoroutines int) *WorkerPool {
	p := WorkerPool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return &p
}

func (p *WorkerPool) Run(w Worker) {
	p.work <- w
}

func (p *WorkerPool) Shutdown () {
	close(p.work)
	p.wg.Wait()
}