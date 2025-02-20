package queue

import (
	"go-parallel_queue/internal/processing/task"
	"sync"
)

type Queue struct {
	mu    sync.RWMutex
	tasks []*task.Task
}

func NewQueue() *Queue {
	return &Queue{
		tasks: make([]*task.Task, 0),
	}
}

func (q *Queue) Append(t *task.Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tasks = append(q.tasks, t)

	return nil
}

func (q *Queue) ShiftUnique(tasks map[string]bool) (*task.Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	emptyTask := task.Task{}
	if len(q.tasks) == 0 {
		return &emptyTask, false
	}

	for i, t := range q.tasks {
		if _, exists := tasks[t.ID]; !exists {
			q.tasks = append(q.tasks[:i], q.tasks[i+1:]...)
			return t, true
		}
	}

	return &emptyTask, false
}

func (q *Queue) State() []*task.Task {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.tasks
}

func (q *Queue) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.tasks)
}
