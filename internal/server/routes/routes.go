package routes

import (
	"go-parallel_queue/internal/server/resolvers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, r *resolvers.Resolver) {
	app.Get("/", r.HomeResolver)
	app.Get("/tasks", r.TasksResolver)
	app.Post("/plan", r.PlanResolver)
}
