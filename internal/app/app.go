package app

import (
	"go-parallel_queue/internal/messages"
	"go-parallel_queue/internal/server/resolvers"
	"go-parallel_queue/internal/server/routes"
	"go-parallel_queue/internal/server/services"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

const FiberShutdownTimeout = 5 * time.Second
const ExecutorShutdownTimeout = 5 * time.Second

func Run() error {
	log.SetFlags(log.Ltime)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./public")

	execService := services.NewExecutorService(nil)
	r := resolvers.NewResolver(execService)

	routes.SetupRoutes(app, r)

	serverShutdown := make(chan struct{})

	go func() {
		<-sigChan
		log.Println(messages.StartShutdownFiber)
		app.ShutdownWithTimeout(FiberShutdownTimeout)
		serverShutdown <- struct{}{}
	}()

	log.Println(app.Stack())

	err := app.Listen(":8080")

	<-serverShutdown
	log.Println(messages.EndShutdownFiber)

	log.Println(messages.StartShutdownExecutorService)

	serviceRet := execService.ShutdownWithTimeout(ExecutorShutdownTimeout)
	if serviceRet {
		log.Println(messages.EndShutdownExecutorServiceSuccess)
	} else {
		log.Println(messages.EndShutdownExecutorServiceTimeout)
	}

	return err
}
