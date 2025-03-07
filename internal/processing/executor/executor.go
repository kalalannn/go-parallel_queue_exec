package executor

import (
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/processing/queue"
	"go-parallel_queue/internal/processing/task"
	"go-parallel_queue/pkg/utils"
	"log"
	"sync"
	"sync/atomic"
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
	CountProcessed int64
	ShutdownFlag   atomic.Bool
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

// return copy of e.activeTasks
func (e *Executor) ActiveTasks() map[string]int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return utils.CopyMap(e.activeTasks)
}

func (e *Executor) Notify() {
	e.cond.Signal()
}

func (e *Executor) Shutdown() {
	e.ShutdownFlag.Store(true)

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
			if e.ShutdownFlag.Load() {
				e.cond.L.Unlock()
				break OuterLoop
			}

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
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.queue.ShiftUnique(e.activeTasks)
}

func (e *Executor) executeTask(t *task.Task) {
	defer func() {
		<-e.workerChan

		// increment countProcessed (++)
		atomic.AddInt64(&e.CountProcessed, 1)

		e.mu.Lock()
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
