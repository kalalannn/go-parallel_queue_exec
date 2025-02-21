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

const FiberShutdownTimeout = 1 * time.Second
const ExecutorShutdownTimeout = 5 * time.Second

type App struct {
	app         *fiber.App
	execService *services.ExecutorService
	sigChan     chan os.Signal
	withHTML    bool
	withWs      bool
}

type AppOptions struct {
	WithHTML bool
	WithWs   bool
}

func NewApp(opts *AppOptions) *App {
	log.SetFlags(log.Ltime)

	if opts == nil {
		opts = &AppOptions{
			WithHTML: true,
			WithWs:   true,
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	return &App{
		sigChan:  sigChan,
		withHTML: opts.WithHTML,
		withWs:   opts.WithWs,
	}
}

func (a *App) SetupApp() {
	if a.withHTML && a.withWs {
		a.setupHTMLWsApp()
	} else if a.withHTML {
		a.setupHTMLApp()
	} else if a.withWs {
		a.setupWsApp()
	} else {
		a.setupRESTApp()
	}
}

func (a *App) setupHTMLWsApp() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./public")

	a.execService = services.NewExecutorService(&services.ExecutorServiceOptions{UseWs: true})

	routes := routes.NewRoutes(app, resolvers.NewResolver(a.execService))
	routes.RESTRoutes()
	routes.HTMLRoutes()
	routes.WSRoutes()

	a.app = app
}

func (a *App) setupHTMLApp() {
	app := fiber.New()

	a.execService = services.NewExecutorService(nil)

	routes := routes.NewRoutes(app, resolvers.NewResolver(a.execService))
	routes.RESTRoutes()
	routes.HTMLRoutes()

	a.app = app
}

func (a *App) setupWsApp() {
	app := fiber.New()

	a.execService = services.NewExecutorService(&services.ExecutorServiceOptions{UseWs: true})

	routes := routes.NewRoutes(app, resolvers.NewResolver(a.execService))
	routes.RESTRoutes()
	routes.WSRoutes()

	a.app = app
}

func (a *App) setupRESTApp() {
	app := fiber.New()

	a.execService = services.NewExecutorService(nil)

	routes := routes.NewRoutes(app, resolvers.NewResolver(a.execService))
	routes.RESTRoutes()

	a.app = app
}

func (a *App) Run() error {
	serverShutdown := make(chan struct{})

	go func() {
		<-a.sigChan
		a.shutdownFiber()
		serverShutdown <- struct{}{}
	}()

	err := a.app.Listen(":8080")

	<-serverShutdown
	a.cleanup()

	return err
}

func (a *App) shutdownFiber() {
	log.Println(messages.StartShutdownServer)

	// shutdown broadcast if for WebSocket
	if a.withWs {
		a.execService.ShutdownBroadcast()
	}

	log.Println(messages.StartShutdownFiber)

	// shutdown fiber
	a.app.ShutdownWithTimeout(FiberShutdownTimeout)
}

func (a *App) cleanup() {
	log.Println(messages.EndShutdownFiber)
	log.Println(messages.StartShutdownExecutorService)

	// shutdown executor
	serviceRet := a.execService.ShutdownWithTimeout(ExecutorShutdownTimeout)

	if serviceRet {
		log.Println(messages.EndShutdownExecutorServiceSuccess)
	} else {
		log.Println(messages.EndShutdownExecutorServiceTimeout)
	}
}
