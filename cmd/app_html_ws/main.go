package main

import (
	"go-parallel_queue/internal/app"
	"log"
)

func main() {
	a := app.NewApp(nil)
	a.SetupApp()
	if err := a.Run(); err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
