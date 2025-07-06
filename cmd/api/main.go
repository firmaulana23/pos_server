package main

import (
	"log"
	"pos-system/internal/app"
	"pos-system/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and start the application
	application := app.NewApp(cfg)
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
