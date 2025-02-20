package executor

import (
	"go-parallel_queue/internal/processing/queue"
	"go-parallel_queue/internal/processing/task"
	"log"
	"sync"
	"time"
)

type Executor struct {
	queue          *queue.Queue
	activeTasks    map[string]bool
	workerChan     chan struct{}
	mu             sync.RWMutex
	wg             sync.WaitGroup
	cond           *sync.Cond
	CountProcessed int
	ShutdownFlag   bool // mimimize padding
}

type ExecutorOptions struct {
	WorkersLimit int
}

func NewExecutor(queue *queue.Queue, executorOptions *ExecutorOptions) *Executor {
	if executorOptions == nil {
		executorOptions = &ExecutorOptions{
			WorkersLimit: 5,
		}
	}
	return &Executor{
		queue:       queue,
		activeTasks: make(map[string]bool),
		cond:        sync.NewCond(&sync.Mutex{}),
		workerChan:  make(chan struct{}, executorOptions.WorkersLimit),
	}
}

func (e *Executor) PlanTasks(tasks ...*task.Task) {
	e.queue.Append(tasks...)
}

func (e *Executor) ActiveTasks() map[string]bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.activeTasks
}

func (e *Executor) PlannedTasks() []*task.Task {
	return e.queue.Tasks()
}

func (e *Executor) Notify() {
	e.cond.Signal()
}

func (e *Executor) Shutdown() {
	e.mu.Lock()
	e.ShutdownFlag = true
	e.mu.Unlock()

	e.Notify()
}

func (e *Executor) Execute(callerWg *sync.WaitGroup) {
OuterLoop:
	for {
		var task *task.Task
		var ok bool

		e.cond.L.Lock()
		for {
			// check shutdown
			e.mu.RLock()
			if e.ShutdownFlag {
				e.mu.RUnlock()
				e.cond.L.Unlock()
				break OuterLoop
			}
			e.mu.RUnlock()

			// next
			task, ok = e.nextTask()
			if ok {
				break
			}

			// wait for notify
			e.cond.Wait()
		}
		e.cond.L.Unlock()

		e.workerChan <- struct{}{}

		e.mu.Lock()
		e.activeTasks[task.ID] = true
		e.mu.Unlock()

		e.wg.Add(1)
		go e.executeTask(task)
	}

	log.Println("Waiting for workers to finish tasks...")
	e.wg.Wait()
	log.Printf("Workers: done. Count processed: %d\n", e.CountProcessed)

	log.Println("Executor: done. Shutdown")
	callerWg.Done()
}

func (e *Executor) nextTask() (*task.Task, bool) {
	return e.queue.ShiftUnique(e.ActiveTasks())
}

func (e *Executor) executeTask(t *task.Task) {
	defer func() {
		<-e.workerChan

		e.mu.Lock()
		e.CountProcessed++
		delete(e.activeTasks, t.ID)
		e.mu.Unlock()

		if e.queue.Len() != 0 {
			e.Notify()
		}

		e.wg.Done()
	}()

	log.Printf("START: %s\n", t.ID)
	time.Sleep(time.Duration(t.Duration) * time.Millisecond)
	log.Printf("DONE: %s\n", t.ID)
}
