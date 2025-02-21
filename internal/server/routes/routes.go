package routes

import (
	"go-parallel_queue/internal/server/resolvers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Routes struct {
	app      *fiber.App
	resolver *resolvers.Resolver
}

func NewRoutes(app *fiber.App, resolver *resolvers.Resolver) *Routes {
	return &Routes{
		app:      app,
		resolver: resolver,
	}
}

func (r *Routes) RESTRoutes() {
	r.app.Get("/tasks", r.resolver.TasksResolver)
	r.app.Post("/plan", r.resolver.PlanResolver)
}

func (r *Routes) HTMLRoutes() {
	r.resolver.UseHTML = true
	r.app.Get("/", r.resolver.HomeResolver)
}

func (r *Routes) WSRoutes() {
	r.resolver.UseWs = true
	r.app.Get("/ws", websocket.New(r.resolver.WebSocketResolver))
}
