package queue_test

import (
	"go-parallel_queue/internal/processing/queue"
	"go-parallel_queue/internal/processing/task"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	q := queue.NewQueue()
	emptyTask := &task.Task{}

	// empty queue
	empty, notOk := q.ShiftUnique(map[string]bool{})

	assert.Zero(t, q.Len())
	assert.False(t, notOk)
	assert.Equal(t, empty, emptyTask)

	// arrange
	task1, task2 := task.NewTask("1", 1), task.NewTask("2", 2)
	q.Append(task1)
	q.Append(task2)

	// shift works
	taskShift, ok := q.ShiftUnique(map[string]bool{"some": true})
	assert.True(t, ok)
	assert.Equal(t, taskShift, task1)
	assert.Equal(t, 1, q.Len())

	// not unique
	taskShift, ok = q.ShiftUnique(map[string]bool{task2.ID: true})
	assert.False(t, ok)
	assert.Equal(t, taskShift, emptyTask)

	// state correct
	order := q.State()
	assert.Equal(t, order, []*task.Task{task2})
}
