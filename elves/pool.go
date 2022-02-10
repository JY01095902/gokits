/*
	Idea of package name "elves" is from "house-elves"
*/

package elves

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
		w := worker{
			isWorking: false,
			locker:    &sync.RWMutex{},
		}
		p.workers = append(p.workers, w)
		w.do(p.tasks)
		p.tasks <- task

		p.locker.Unlock()
		return
	}
	p.locker.Unlock()

retrieve:
	for i := range p.workers {
		w := p.workers[i]
		w.locker.RLock()
		isWorking := w.isWorking
		w.locker.RUnlock()
		if !isWorking {
			w.do(p.tasks)
			p.tasks <- task

			return
		}
	}

	goto retrieve
}

func (p *Pool) Destroy() {
	close(p.tasks)
}
