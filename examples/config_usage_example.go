package main

import (
	"context"
	"fmt"
	"log"

	"lam-phuong-api/internal/config"
)

// Example 1: Using Global Config (Simple)
func exampleGlobalConfig() {
	// Access config anywhere using Get()
	cfg := config.Get()
	
	fmt.Printf("Server Port: %s\n", cfg.Server.Port)
	fmt.Printf("Airtable Base ID: %s\n", cfg.Airtable.BaseID)
}

// Example 2: Using Config as Parameter (Dependency Injection - Recommended)
func exampleWithParameter(cfg *config.Config) {
	// Use the passed config
	fmt.Printf("Server running on: %s\n", cfg.ServerAddress())
	
	// Create Airtable client using config
	client, err := cfg.NewAirtableClient()
	if err != nil {
		log.Printf("Error creating client: %v", err)
		return
	}
	
	// Use the client
	records, err := client.ListRecords(context.Background(), "Books", nil)
	if err != nil {
		log.Printf("Error listing records: %v", err)
		return
	}
	
	fmt.Printf("Found %d records\n", len(records))
}

// Example 3: Service with Config Dependency
type BookService struct {
	cfg *config.Config
}

func NewBookService(cfg *config.Config) *BookService {
	return &BookService{cfg: cfg}
}

func (s *BookService) GetAirtableConfig() (string, string) {
	return s.cfg.Airtable.APIKey, s.cfg.Airtable.BaseID
}

func (s *BookService) GetServerInfo() string {
	return s.cfg.ServerAddress()
}

// Example 4: Handler with Config
type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

func (h *ConfigHandler) GetConfigInfo() map[string]interface{} {
	return map[string]interface{}{
		"server": map[string]interface{}{
			"port": h.cfg.Server.Port,
			"host": h.cfg.Server.Host,
		},
		"airtable": map[string]interface{}{
			"base_id": h.cfg.Airtable.BaseID,
			// Don't expose API key in responses!
		},
	}
}

func main() {
	// Load config once at startup
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Example 1: Global config
	fmt.Println("=== Example 1: Global Config ===")
	exampleGlobalConfig()

	// Example 2: Pass config as parameter
	fmt.Println("\n=== Example 2: Config as Parameter ===")
	exampleWithParameter(cfg)

	// Example 3: Service with config
	fmt.Println("\n=== Example 3: Service with Config ===")
	bookService := NewBookService(cfg)
	fmt.Printf("Server info: %s\n", bookService.GetServerInfo())

	// Example 4: Handler with config
	fmt.Println("\n=== Example 4: Handler with Config ===")
	configHandler := NewConfigHandler(cfg)
	info := configHandler.GetConfigInfo()
	fmt.Printf("Config info: %+v\n", info)
}

