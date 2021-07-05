package elves

type worker struct {
	isWorking bool
}

func (w *worker) do(tasks <-chan func()) {
	w.isWorking = true
	go func() {
		defer func() {
			w.isWorking = false
		}()

		task := <-tasks
		if task == nil {
			return
		}

		task()
	}()
}
