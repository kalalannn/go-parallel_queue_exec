package executor_test

import (
	"go-parallel_queue/internal/processing/executor"
	"go-parallel_queue/internal/processing/task"
	"go-parallel_queue/pkg/utils"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var sleepTime, gap = 75, 3

func TestEmpty(t *testing.T) {
	// arrange
	e := executor.NewExecutor(nil)

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	// act
	e.Shutdown()
	ourWg.Wait()

	// assert
	assert.Equal(t, 0, e.CountProcessed)
	assert.True(t, e.ShutdownFlag)
}

func TestShutdownFasterThanOne(t *testing.T) {
	// arrange
	e := executor.NewExecutor(nil)

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	e.ScheduleTasks(task.NewTask("1", 1))

	// act
	e.Notify()
	e.Shutdown()
	ourWg.Wait()

	// assert
	assert.Equal(t, 0, e.CountProcessed)
	assert.True(t, e.ShutdownFlag)
}

func TestOneProcessed(t *testing.T) {
	// arrange
	e := executor.NewExecutor(nil)

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	e.ScheduleTasks(task.NewTask("1", 1))

	// act
	e.Notify()
	time.Sleep(20 * time.Millisecond)

	e.Shutdown()
	ourWg.Wait()

	// assert
	assert.Equal(t, 1, e.CountProcessed)
	assert.True(t, e.ShutdownFlag)
}

func TestNobodyBlocked(t *testing.T) {
	// arrange
	e := executor.NewExecutor(&executor.ExecutorOptions{WorkersLimit: 1})

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	e.ScheduleTasks(task.NewTask("1", 15))

	// act
	e.Notify()
	time.Sleep(10 * time.Millisecond)

	e.Shutdown()
	ourWg.Wait()

	// assert
	assert.Equal(t, 1, e.CountProcessed)
	assert.True(t, e.ShutdownFlag)
}

// if failing increase gap or sleepTime
func TestWorkersCount3(t *testing.T) {
	// arrange
	e := executor.NewExecutor(&executor.ExecutorOptions{WorkersLimit: 3})

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	e.ScheduleTasks(
		task.NewTask("1", sleepTime*gap),
		task.NewTask("2", sleepTime*gap),
		task.NewTask("3", sleepTime*gap),
		task.NewTask("4", 1),
		task.NewTask("5", 1),
	)

	// act
	e.Notify()

	// wait
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	// assert
	activeTasks := utils.MapKeys(e.ActiveTasks())
	assert.Equal(t, 3, len(activeTasks))

	// cleanup
	e.Shutdown()
	ourWg.Wait()
}

func TestUniqueExecution(t *testing.T) {
	// arrange
	e := executor.NewExecutor(&executor.ExecutorOptions{WorkersLimit: 2})

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	e.ScheduleTasks(
		task.NewTask("1", sleepTime*gap),
		task.NewTask("1", 1),
	)

	// act
	e.Notify()

	// wait
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	// take stateshot
	activeTasks := utils.MapKeys(e.ActiveTasks())
	assert.Equal(t, 1, len(activeTasks))

	// wait
	time.Sleep(time.Duration(sleepTime*gap*2) * time.Millisecond)

	e.Shutdown()
	ourWg.Wait()

	// assert
	assert.Equal(t, 2, e.CountProcessed)
}
