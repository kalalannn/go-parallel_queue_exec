package queue_test

import (
	"go-parallel_queue/internal/processing/queue"
	"go-parallel_queue/internal/processing/task"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	q := queue.NewQueue()

	// empty queue
	assert.Zero(t, q.Len())

	empty, notOk := q.Shift()
	assert.False(t, notOk)
	assert.Equal(t, empty, task.Task{})

	task1, task2 := task.NewTask("1", 1), task.NewTask("2", 2)
	q.Append(task1)
	q.Append(task2)

	// already exists
	err := q.Append(task1)
	assert.Error(t, err)

	// shift works
	taskShift, ok := q.Shift()
	assert.True(t, ok)
	assert.Equal(t, taskShift, task1)
	assert.Equal(t, 1, q.Len())

	// state correct
	tasks, order := q.State()
	assert.Equal(t, tasks, map[string]task.Task{"2": task2})
	assert.Equal(t, order, []string{"2"})
}
