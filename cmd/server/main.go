package main

import (
	"log"

	"lam-phuong-api/internal/book"
	"lam-phuong-api/internal/config"
	"lam-phuong-api/internal/location"
	"lam-phuong-api/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize seed data
	bookSeed := []book.Book{
		{ID: "1", Title: "The Go Programming Language", Author: "Alan A. A. Donovan"},
		{ID: "2", Title: "Introducing Go", Author: "Caleb Doxsey"},
	}

	locationSeed := []location.Location{
		{ID: "1", Name: "Main Library", Address: "123 Main St", City: "Go City"},
		{ID: "2", Name: "West Branch", Address: "456 Elm St", City: "Go City"},
	}

	bookRepo := book.NewInMemoryRepository(bookSeed)
	bookHandler := book.NewHandler(bookRepo)

	locationRepo := location.NewInMemoryRepository(locationSeed)
	locationHandler := location.NewHandler(locationRepo)

	router := server.NewRouter(bookHandler, locationHandler)

	// Use server address from config
	serverAddr := cfg.ServerAddress()
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
