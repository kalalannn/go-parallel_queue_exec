package executor

import (
	"fmt"
	"go-parallel_queue/internal/processing/queue"
	"go-parallel_queue/internal/processing/task"
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

func (e *Executor) Execute(callerWg *sync.WaitGroup) {
	for {
		// wait for notify
		e.cond.L.Lock()
		for e.queue.Len() == 0 && !e.ShutdownFlag {
			e.cond.Wait()
		}
		e.cond.L.Unlock()

		e.workerChan <- struct{}{}

		e.mu.RLock()
		if e.ShutdownFlag {
			e.mu.RUnlock()
			break
		}
		e.mu.RUnlock()

		task, ok := e.queue.ShiftUnique(e.State())
		if !ok {
			// should be always ok, but just in case
			<-e.workerChan
			continue
		}

		e.mu.Lock()
		e.activeTasks[task.ID] = true
		e.mu.Unlock()

		e.wg.Add(1)
		go e.executeTask(task)
	}

	fmt.Println("Waiting for workers to finish tasks...")
	e.wg.Wait()
	fmt.Printf("Workers: done. Count processed: %d\n", e.CountProcessed)

	fmt.Println("Executor: done. Shutdown")
	callerWg.Done()
}

func (e *Executor) executeTask(t task.Task) {
	defer func() {
		<-e.workerChan

		e.mu.Lock()
		e.CountProcessed++
		delete(e.activeTasks, t.ID)
		e.mu.Unlock()

		e.wg.Done()
	}()

	fmt.Printf("Task: %s\n", t.ID)
	time.Sleep(time.Duration(t.Duration) * time.Millisecond)
	fmt.Printf("Done: %s done\n", t.ID)
}

func (e *Executor) Notify() {
	e.cond.Signal()
}

func (e *Executor) Shutdown() {
	e.mu.Lock()
	e.ShutdownFlag = true
	e.mu.Unlock()

	// if waiting for new tasks
	e.Notify()

	<-e.workerChan
}

func (e *Executor) State() map[string]bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.activeTasks
}
