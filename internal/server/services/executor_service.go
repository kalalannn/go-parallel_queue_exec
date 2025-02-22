package services

import (
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/processing/executor"
	"go-parallel_queue/internal/processing/task"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

const (
	defaultExecutorWorkersLimit = 5
	defaultUseWebSockets        = false
)

type ExecutorService struct {
	exec          *executor.Executor
	ourWg         *sync.WaitGroup
	stopBroadcast chan struct{}
	broadcastDone chan struct{}
	wsMu          sync.RWMutex
	wsClients     map[*websocket.Conn]bool
	wsUpdateChan  chan map[string]any
}

type ExecutorServiceOptions struct {
	WorkersLimit int
	UseWs        bool
}

func NewExecutorService(opts *ExecutorServiceOptions) *ExecutorService {
	if opts == nil {
		opts = &ExecutorServiceOptions{
			WorkersLimit: defaultExecutorWorkersLimit,
			UseWs:        defaultUseWebSockets,
		}
	} else if opts.WorkersLimit == 0 {
		opts.WorkersLimit = defaultExecutorWorkersLimit
	}

	ourWg := sync.WaitGroup{}

	var exec *executor.Executor
	var execService ExecutorService
	if opts.UseWs {
		wsUpdateChan := make(chan map[string]any)
		exec = executor.NewExecutor(&executor.ExecutorOptions{
			WorkersLimit: opts.WorkersLimit,
			UpdatesChan:  wsUpdateChan,
		})

		execService = ExecutorService{
			exec:          exec,
			ourWg:         &ourWg,
			wsClients:     make(map[*websocket.Conn]bool),
			wsUpdateChan:  wsUpdateChan,
			stopBroadcast: make(chan struct{}),
			broadcastDone: make(chan struct{}),
		}

		ourWg.Add(1)
		go execService.BroadcastUpdates()

	} else {
		exec = executor.NewExecutor(&executor.ExecutorOptions{
			WorkersLimit: opts.WorkersLimit,
		})

		execService = ExecutorService{
			exec:  exec,
			ourWg: &ourWg,
		}
	}

	ourWg.Add(1)
	go exec.Execute(&ourWg)

	return &execService
}

func (s *ExecutorService) AddWebSocketClient(client *websocket.Conn) {
	s.wsMu.Lock()
	defer s.wsMu.Unlock()
	s.wsClients[client] = true
	log.Println(messages.NewWebSocketClient)
}

func (s *ExecutorService) RemoveWebSocketClient(client *websocket.Conn) {
	s.wsMu.Lock()
	defer s.wsMu.Unlock()
	delete(s.wsClients, client)
	log.Println(messages.RemoveWebSocketClient)
}

func (s *ExecutorService) WebSocketClientsCount() int {
	s.wsMu.RLock()
	defer s.wsMu.RUnlock()
	return len(s.wsClients)
}

func (s *ExecutorService) BroadcastUpdates() {
	for {
		select {
		case update := <-s.wsUpdateChan:
			s.wsMu.RLock()
			for client := range s.wsClients {
				if err := client.WriteJSON(update); err != nil {
					log.Printf(messages.BroadcastFailed, err)
					client.Close()
					s.RemoveWebSocketClient(client)
				}
			}
			s.wsMu.RUnlock()
		case <-s.stopBroadcast:
			log.Println(messages.StopBroadcastUpdatesStart)
			for client := range s.wsClients {
				client.Close()
				s.RemoveWebSocketClient(client)
			}
			s.ourWg.Done()
			s.broadcastDone <- struct{}{}
			return
		}
	}
}

func (s *ExecutorService) ShutdownBroadcast() {
	s.stopBroadcast <- struct{}{}
	<-s.broadcastDone
	log.Println(messages.StopBroadcastUpdatesEnd)
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
	count := len(data)
	tasks := make([]*task.Task, count)
	i := 0
	for k, v := range data {
		tasks[i] = task.NewTask(k, v)
		i++
	}
	s.exec.PlanTasks(tasks...)
	s.exec.Notify()
}
