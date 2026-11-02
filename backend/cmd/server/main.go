package main

import (
	"log"
	"net/http"

	"github.com/kwizera/university-invoices/internal/config"
	"github.com/kwizera/university-invoices/internal/routes"
)

func main() {
	log.Println("Starting University of Kigali Invoices API Server...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("Database: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Initialize database connection
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Setup routes
	router := routes.SetupRouter(db)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	log.Println("Available endpoints:")
	log.Println("  GET    /health")
	log.Println("  POST   /api/auth/login")
	log.Println("  POST   /api/auth/register")
	log.Println("  GET    /api/users")
	log.Println("  GET    /api/students")
	log.Println("  GET    /api/fees")
	log.Println("  GET    /api/invoices")
	log.Println("  GET    /api/payments")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
