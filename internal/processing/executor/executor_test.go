package executor_test

import (
	"go-parallel_queue/internal/processing/executor"
	"go-parallel_queue/internal/processing/queue"
	"go-parallel_queue/internal/processing/task"
	"go-parallel_queue/pkg/utils"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestEmpty(t *testing.T) {
	// arrange
	q := queue.NewQueue()
	e := executor.NewExecutor(q, nil)

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
	q := queue.NewQueue()
	e := executor.NewExecutor(q, nil)

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	q.Append(task.NewTask("1", 1))

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
	q := queue.NewQueue()
	e := executor.NewExecutor(q, nil)

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	q.Append(task.NewTask("1", 1))

	// act
	e.Notify()
	time.Sleep(10 * time.Millisecond)

	e.Shutdown()
	ourWg.Wait()

	// assert
	assert.Equal(t, 1, e.CountProcessed)
	assert.True(t, e.ShutdownFlag)
}

func TestNobodyBlocked(t *testing.T) {
	// arrange
	q := queue.NewQueue()
	e := executor.NewExecutor(q, &executor.ExecutorOptions{WorkersLimit: 1})

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	q.Append(task.NewTask("1", 15))

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
	q := queue.NewQueue()
	e := executor.NewExecutor(q, &executor.ExecutorOptions{WorkersLimit: 3})

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go e.Execute(&ourWg)

	sleepTime := 75 // ms
	gap := 3

	q.Append(task.NewTask("1", sleepTime*gap))
	q.Append(task.NewTask("2", sleepTime*2))
	q.Append(task.NewTask("3", sleepTime*gap/2))
	q.Append(task.NewTask("4", 1))
	q.Append(task.NewTask("5", 1))

	// act
	e.Notify()
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	activeTasks := utils.MapKeys(e.State())

	e.Shutdown()
	ourWg.Wait()

	// assert
	assert.Equal(t, 3, len(activeTasks))
}
