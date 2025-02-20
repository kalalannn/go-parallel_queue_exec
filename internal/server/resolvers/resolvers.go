package resolvers

import (
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/server/services"

	"github.com/gofiber/fiber/v2"
)

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
