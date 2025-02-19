package routes

import (
	"go-parallel_queue/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/status", handlers.StatusHandler)
	app.Post("/plan", handlers.PlanHandler)
}
