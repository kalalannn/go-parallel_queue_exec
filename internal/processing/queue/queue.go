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

func (q *Queue) Append(tasks ...*task.Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tasks = append(q.tasks, tasks...)

	return nil
}

func (q *Queue) ShiftUnique(excludeTasks map[string]int) (*task.Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) == 0 {
		return nil, false
	}

	for i, t := range q.tasks {
		if _, exists := excludeTasks[t.ID]; !exists {
			q.tasks = append(q.tasks[:i], q.tasks[i+1:]...)
			return t, true
		}
	}

	return nil, false
}

func (q *Queue) Tasks() []*task.Task {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.tasks
}

func (q *Queue) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.tasks)
}
