package main

import (
	"go-parallel_queue/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
