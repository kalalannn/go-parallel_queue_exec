package main

import (
	"go-parallel_queue/internal/app"
	"go-parallel_queue/pkg/utils"
	"log"
)

func main() {
	// setup logger
	log.SetFlags(log.Ltime)
	// load config
	config := utils.MustLoadConfig()

	// create && setup app
	a := app.NewApp(config, &app.AppOptions{WithHTML: false, WithWs: false})
	a.SetupApp()

	// run
	if err := a.Run(); err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
