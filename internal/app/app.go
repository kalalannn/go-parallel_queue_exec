package app

import (
	"fmt"
	"go-parallel_queue/internal/config"
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

type App struct {
	app         *fiber.App
	config      *config.Config
	execService *services.ExecutorService
	sigChan     chan os.Signal
	withHTML    bool
	withWs      bool
}

type AppOptions struct {
	WithHTML bool
	WithWs   bool
}

func NewApp(config *config.Config, opts *AppOptions) *App {
	if config == nil {
		log.Fatalln("no config for app")
	}

	if opts == nil {
		opts = &AppOptions{
			WithHTML: true,
			WithWs:   true,
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	return &App{
		config:   config,
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

func (a *App) fiberAppWithHTML() *fiber.App {
	engine := html.New(a.config.App.ViewsFolder, a.config.App.TemplatesExt)
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static(a.config.App.StaticEndpoint, a.config.App.PublicFolder)

	return app
}

func (a *App) newExecService(useWs bool) *services.ExecutorService {
	return services.NewExecutorService(&services.ExecutorServiceOptions{
		WorkersLimit: a.config.ExecutorService.WorkersLimit,
		UseWs:        useWs,
	})
}

func (a *App) setupHTMLWsApp() {
	app := a.fiberAppWithHTML()

	a.execService = a.newExecService(true)

	routes := routes.NewRoutes(app, resolvers.NewResolver(a.execService))
	routes.RESTRoutes()
	routes.HTMLRoutes()
	routes.WSRoutes()

	a.app = app
}

func (a *App) setupHTMLApp() {
	app := a.fiberAppWithHTML()

	a.execService = a.newExecService(false)

	routes := routes.NewRoutes(app, resolvers.NewResolver(a.execService))
	routes.RESTRoutes()
	routes.HTMLRoutes()

	a.app = app
}

func (a *App) setupWsApp() {
	app := fiber.New()

	a.execService = a.newExecService(true)

	routes := routes.NewRoutes(app, resolvers.NewResolver(a.execService))
	routes.RESTRoutes()
	routes.WSRoutes()

	a.app = app
}

func (a *App) setupRESTApp() {
	app := fiber.New()

	a.execService = a.newExecService(false)

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

	err := a.app.Listen(fmt.Sprintf("%s:%d", a.config.App.Host, a.config.App.Port))

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
	a.app.ShutdownWithTimeout(time.Duration(a.config.App.FiberShutdownTimeoutMs) * time.Millisecond)
}

func (a *App) cleanup() {
	log.Println(messages.EndShutdownFiber)
	log.Println(messages.StartShutdownExecutorService)

	// shutdown executor
	serviceRet := a.execService.ShutdownWithTimeout(time.Duration(a.config.App.ExecutorShutdownTimeoutMs) * time.Millisecond)

	if serviceRet {
		log.Println(messages.EndShutdownExecutorServiceSuccess)
	} else {
		log.Println(messages.EndShutdownExecutorServiceTimeout)
	}
}
