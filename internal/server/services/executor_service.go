package services

import (
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/processing/executor"
	"go-parallel_queue/internal/processing/task"
	"log"
	"sync"
	"time"
)

type ExecutorService struct {
	exec  *executor.Executor
	ourWg *sync.WaitGroup
}

type ExecutorServiceOptions struct {
	WorkersLimit int
}

func NewExecutorService(opts *ExecutorServiceOptions) *ExecutorService {
	var exec *executor.Executor
	if opts == nil {
		exec = executor.NewExecutor(nil)
	} else {
		exec = executor.NewExecutor(&executor.ExecutorOptions{WorkersLimit: opts.WorkersLimit})
	}
	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go exec.Execute(&ourWg)
	return &ExecutorService{
		exec:  exec,
		ourWg: &ourWg,
	}
}

func (s *ExecutorService) Shutdown() bool {
	s.exec.Shutdown()
	s.ourWg.Wait()
	return true
}

func (s *ExecutorService) ShutdownWithTimeout(timeout time.Duration) bool {
	s.exec.Shutdown()
	log.Printf(messages.WaitForExecutorShutdownWithTimeout, timeout.String())

	c := make(chan struct{})
	go func() {
		defer close(c)
		s.ourWg.Wait()
	}()
	select {
	case <-c:
		log.Println(messages.ExecutorShutdownFinished)
		return true
	case <-time.After(timeout):
		log.Println(messages.ExecutorShutdownTimeout)
		return false
	}
}

func (s *ExecutorService) ActiveTasks() map[string]int {
	return s.exec.ActiveTasks()
}

func (s *ExecutorService) PlannedTasks() []map[string]int {
	plannedTasks := make([]map[string]int, 0)
	for _, t := range s.exec.PlannedTasks() {
		plannedTasks = append(plannedTasks, map[string]int{t.ID: t.Duration})
	}
	return plannedTasks
}

func (s *ExecutorService) PlanExecuteTasks(data map[string]int) {
	for k, v := range data {
		s.exec.PlanTasks(task.NewTask(k, v))
	}
	s.exec.Notify()
}
