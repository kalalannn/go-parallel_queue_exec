package resolvers

import (
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/server/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

const restAccepts = "application/json"

const (
	wsMessageTag   = "message"
	restMessageTag = "message"
	restErrorTag   = "error"
)
const (
	tasksActiveTag    = "active"
	tasksScheduledTag = "scheduled"
)
const (
	wsReadErrorMsg  = "Read error: "
	wsWriteErrorMsg = "Write error: "
)

type Resolver struct {
	execService *services.ExecutorService
	UseWs       bool
	UseHTML     bool
}

func NewResolver(execService *services.ExecutorService) *Resolver {
	return &Resolver{
		execService: execService,
	}
}

func (r *Resolver) HomeResolver(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"WithWs": r.UseWs,
	})
}

func (r *Resolver) WebSocketResolver(c *websocket.Conn) {
	defer c.Close()

	r.execService.AddWebSocketClient(c)
	defer r.execService.RemoveWebSocketClient(c)

	_, _, err := c.ReadMessage()
	if err != nil {
		log.Println(wsReadErrorMsg, err)
		return
	}

	if err := c.WriteJSON(map[string]string{wsMessageTag: messages.WelcomeMessage}); err != nil {
		log.Println(wsWriteErrorMsg, err)
		return
	}

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println(wsReadErrorMsg, err)
			return
		}

		if err := c.WriteJSON(map[string]string{wsMessageTag: messages.UseRESTMessage}); err != nil {
			log.Println(wsWriteErrorMsg, err)
			return
		}
	}
}

func (r *Resolver) TasksResolver(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		tasksActiveTag:    r.execService.ActiveTasks(),
		tasksScheduledTag: r.execService.ScheduledTasks(),
	})
}

func (r *Resolver) ScheduleResolver(c *fiber.Ctx) error {
	c.Accepts(restAccepts)
	var data map[string]int
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			restErrorTag: messages.InvalidJSON,
		})
	}

	r.execService.ScheduleExecuteTasks(data)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		restMessageTag: messages.OK,
	})
}
