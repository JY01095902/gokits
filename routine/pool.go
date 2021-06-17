package routine

import (
	"sync"
)

type Task struct {
	Id      string
	Job     func() error
	Message string
}

type TaskResult struct {
	TaskId  string
	Status  string
	Error   error
	Message string
}

type ResultHandler func(result TaskResult)

type Pool struct {
	workersCount   int
	tasksChannel   <-chan Task
	waitGroup      *sync.WaitGroup
	resultsChannel chan TaskResult
}

func NewPool(workersCnt int, tasksChan <-chan Task) Pool {
	pool := Pool{
		workersCount: workersCnt,
		tasksChannel: tasksChan,
		waitGroup:    &sync.WaitGroup{},
	}

	return pool
}

func NewPoolWithResultHandler(workersCnt int, tasksChan <-chan Task, resultHandler ResultHandler) Pool {
	pool := NewPool(workersCnt, tasksChan)
	pool.resultsChannel = make(chan TaskResult)

	go func() {
		for res := range pool.resultsChannel {
			resultHandler(res)
		}
	}()

	return pool
}

func (pool Pool) run() {
	pool.waitGroup.Add(1)

	var wg sync.WaitGroup
	wg.Add(pool.workersCount)
	for i := 1; i <= pool.workersCount; i++ {
		worker := NewWorker(i)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			worker.Do(pool.tasksChannel, pool.resultsChannel)
		}(&wg)
	}
	wg.Wait()

	if pool.resultsChannel != nil {
		close(pool.resultsChannel)
	}

	pool.waitGroup.Done()
}

func (pool Pool) Run() {
	go pool.run()
}

func (pool Pool) Wait() {
	pool.waitGroup.Wait()
}
