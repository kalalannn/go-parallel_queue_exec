package executor

import (
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/processing/queue"
	"go-parallel_queue/internal/processing/task"
	"log"
	"sync"
	"time"
)

const defaultWorkersLimit = 5

type Executor struct {
	queue          *queue.Queue
	activeTasks    map[string]int
	workerChan     chan struct{}
	updatesChan    chan map[string]any
	mu             sync.RWMutex
	wg             sync.WaitGroup
	cond           *sync.Cond
	CountProcessed int
	ShutdownFlag   bool // mimimize padding
}

type ExecutorOptions struct {
	WorkersLimit int
	UpdatesChan  chan map[string]any
}

func NewExecutor(executorOptions *ExecutorOptions) *Executor {
	if executorOptions == nil {
		executorOptions = &ExecutorOptions{
			WorkersLimit: defaultWorkersLimit,
			UpdatesChan:  nil,
		}
	} else if executorOptions.WorkersLimit == 0 {
		executorOptions.WorkersLimit = defaultWorkersLimit
	}
	return &Executor{
		queue:       queue.NewQueue(),
		activeTasks: make(map[string]int),
		cond:        sync.NewCond(&sync.Mutex{}),
		workerChan:  make(chan struct{}, executorOptions.WorkersLimit),
		updatesChan: executorOptions.UpdatesChan,
	}
}

func (e *Executor) ScheduleTasks(tasks ...*task.Task) {
	e.queue.Append(tasks...)

	if e.updatesChan != nil {
		e.updatesChan <- map[string]any{
			messages.ScheduledTag: tasks,
		}
	}
}

func (e *Executor) ScheduledTasks() []*task.Task {
	return e.queue.Tasks()
}

func (e *Executor) ActiveTasks() map[string]int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.activeTasks
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

		if e.updatesChan != nil {
			e.updatesChan <- map[string]any{
				messages.NextTag: task,
			}
		}

		e.workerChan <- struct{}{}

		e.mu.Lock()
		e.activeTasks[task.ID] = task.Duration
		e.mu.Unlock()

		if e.updatesChan != nil {
			e.updatesChan <- map[string]any{
				messages.StartTag: task,
			}
		}

		e.wg.Add(1)
		go e.executeTask(task)
	}

	log.Println(messages.WaitForWorkers)

	e.wg.Wait()

	log.Printf(messages.WorkersDone, e.CountProcessed)
	log.Println(messages.ExecutorDoneShutdown)

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

		if e.updatesChan != nil {
			e.updatesChan <- map[string]any{
				messages.DoneTag: t,
			}
		}

		if e.queue.Len() != 0 {
			e.Notify()
		}

		e.wg.Done()
	}()

	log.Printf(messages.ExecuteStart, t.ID)
	time.Sleep(time.Duration(t.Duration) * time.Millisecond)
	log.Printf(messages.ExecuteDone, t.ID)
}
