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
}

func NewExecutorService(opts *ExecutorServiceOptions) *ExecutorService {
	var exec *executor.Executor
	wsUpdateChan := make(chan map[string]any)

	if opts == nil {
		exec = executor.NewExecutor(&executor.ExecutorOptions{
			UpdatesChan: wsUpdateChan,
		})
	} else {
		exec = executor.NewExecutor(&executor.ExecutorOptions{
			WorkersLimit: opts.WorkersLimit,
			UpdatesChan:  wsUpdateChan,
		})
	}

	ourWg := sync.WaitGroup{}
	ourWg.Add(1)
	go exec.Execute(&ourWg)

	s := ExecutorService{
		exec:          exec,
		ourWg:         &ourWg,
		wsClients:     make(map[*websocket.Conn]bool),
		wsUpdateChan:  wsUpdateChan,
		stopBroadcast: make(chan struct{}),
		broadcastDone: make(chan struct{}),
	}

	ourWg.Add(1)
	go s.BroadcastUpdates()

	return &s
}

func (s *ExecutorService) AddWebSocketClient(client *websocket.Conn) {
	s.wsMu.Lock()
	defer s.wsMu.Unlock()
	s.wsClients[client] = true
	log.Println("Added new WebSocket client")
}

func (s *ExecutorService) RemoveWebSocketClient(client *websocket.Conn) {
	s.wsMu.Lock()
	defer s.wsMu.Unlock()
	delete(s.wsClients, client)
	log.Println("Removed WebSocket client")
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
					log.Printf("Failed to send update to client: %s, disconnecting\n", err)
					client.Close()
					s.RemoveWebSocketClient(client)
				}
			}
			s.wsMu.RUnlock()
		case <-s.stopBroadcast:
			log.Println("Stopping broadcast updates...")
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
	log.Println("Broadcast updates stopped.")
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
