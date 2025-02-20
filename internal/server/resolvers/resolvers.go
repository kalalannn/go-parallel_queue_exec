package resolvers

import (
	"encoding/json"
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/server/services"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

const wsTimeout = 2 * time.Second

type Resolver struct {
	execService *services.ExecutorService
}

func NewResolver(execService *services.ExecutorService) *Resolver {
	return &Resolver{
		execService: execService,
	}
}

func (r *Resolver) HomeResolver(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Name": "Go Parallel Queue",
	})
}

func (r *Resolver) WebSocketResolver(c *websocket.Conn) {
	defer c.Close()

	r.execService.AddWebSocketClient(c)
	defer r.execService.RemoveWebSocketClient(c)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		_ = msg

		response := map[string]any{
			"activeTasks":  r.execService.ActiveTasks(),
			"plannedTasks": r.execService.PlannedTasks(),
		}
		responseMsg, _ := json.Marshal(response)

		if err := c.WriteMessage(websocket.TextMessage, responseMsg); err != nil {
			log.Println("Write error:", err)
			return
		}

		// time.Sleep(wsTimeout * time.Second)
	}
}

func (r *Resolver) TasksResolver(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"active":  r.execService.ActiveTasks(),
		"planned": r.execService.PlannedTasks(),
	})
}

func (r *Resolver) PlanResolver(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var data map[string]int
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": messages.InvalidJSON,
		})
	}

	r.execService.PlanExecuteTasks(data)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": messages.OK,
	})
}
