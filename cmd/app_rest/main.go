package main

import (
	"go-parallel_queue/internal/app"
	"log"
)

func main() {
	a := app.NewApp(&app.AppOptions{WithHTML: false, WithWs: false})
	a.SetupApp()
	if err := a.Run(); err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
