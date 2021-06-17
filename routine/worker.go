package routine

type Worker struct {
	id int
}

func NewWorker(id int) Worker {
	return Worker{
		id: id,
	}
}

func (worker Worker) Do(tasksChan <-chan Task, resultsChan chan<- TaskResult) {
	for task := range tasksChan {
		result := TaskResult{
			TaskId:  task.Id,
			Status:  "SUCCESSFUL",
			Message: task.Message,
		}

		if err := task.Job(); err != nil {
			result.Status = "FAILED"
			result.Error = err
		}

		if resultsChan != nil {
			resultsChan <- result
		}
	}
}
