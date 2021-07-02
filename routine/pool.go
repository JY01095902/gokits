package routine

import (
	"sync"
)

type Pool struct {
	capacity int
	tasks    chan func()
	workers  []worker
	locker   sync.Locker
}

func NewPool(cap int) (Pool, error) {
	if cap <= 0 {
		return Pool{}, ErrInvalidCapacity
	}

	pool := Pool{
		capacity: cap,
		tasks:    make(chan func()),
		workers:  []worker{},
		locker:   &sync.Mutex{},
	}

	return pool, nil
}

func (p *Pool) Execute(task func()) {
	p.locker.Lock()
	if len(p.workers) < int(p.capacity) {
		w := worker{}
		p.workers = append(p.workers, w)
		w.do(p.tasks)
		p.tasks <- task

		p.locker.Unlock()
		return
	}
	p.locker.Unlock()

retrieve:
	for i := range p.workers {
		if !p.workers[i].isWorking {
			p.workers[i].do(p.tasks)
			p.tasks <- task

			return
		}
	}

	goto retrieve
}

func (p *Pool) Destroy() {
	close(p.tasks)
}
