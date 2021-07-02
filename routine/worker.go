package routine

type worker struct {
	isWorking bool
}

func (w *worker) do(tasks <-chan func()) {
	w.isWorking = true
	go func() {
		defer func() {
			w.isWorking = false
		}()

		for task := range tasks {
			if task == nil {
				return
			}

			task()

			if w.isWorking {
				return
			}
		}
	}()
}
