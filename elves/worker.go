package elves

import (
	"sync"
)

type worker struct {
	isWorking bool
	locker    *sync.RWMutex
}

func (w *worker) do(tasks <-chan func()) {
	w.locker.Lock()
	w.isWorking = true
	w.locker.Unlock()

	go func() {
		defer func() {
			w.locker.Lock()
			w.isWorking = false
			w.locker.Unlock()
		}()

		task := <-tasks
		if task == nil {
			return
		}

		task()
	}()
}
