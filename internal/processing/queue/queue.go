package queue

import (
	"fmt"
	"go-parallel_queue/internal/processing/task"
	"sync"
)

type Queue struct {
	mu    sync.RWMutex
	tasks map[string]task.Task
	order []string
}

func NewQueue() *Queue {
	return &Queue{
		tasks: make(map[string]task.Task),
		order: make([]string, 0),
	}
}

func (q *Queue) Append(t task.Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.tasks[t.ID]; exists {
		return fmt.Errorf("task %s already exists", t.ID)
	}

	q.tasks[t.ID] = t
	q.order = append(q.order, t.ID)

	return nil
}

func (q *Queue) Shift() (task.Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.order) == 0 {
		return task.Task{}, false
	}

	id := q.order[0]
	t := q.tasks[id]

	delete(q.tasks, id)
	q.order = q.order[1:]

	return t, true
}

func (q *Queue) State() (map[string]task.Task, []string) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.tasks, q.order
}

func (q *Queue) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.order)
}
