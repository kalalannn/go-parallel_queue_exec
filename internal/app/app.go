package app

import (
	"go-parallel_queue/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Run() error {
	log.SetFlags(log.Ltime)
	app := fiber.New()
	routes.SetupRoutes(app)
	return app.Listen(":8080")
}
