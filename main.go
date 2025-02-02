// /////////////////////////////////////////////////////////////////////
// main.go
// /////////////////////////////////////////////////////////////////////
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"logging_microservice/db"
	"logging_microservice/handlers"
)

func main() {
	log.Println("Starting log service")
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	if err := db.InitDB("data/logs.db"); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.DB.Close()

	// Create tables if not present
	if err := db.CreateTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Register routes with CORS middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/logs", handlers.LogsHandler)
	mux.HandleFunc("/logs/", handlers.LogsHandler)
	mux.HandleFunc("/loggers", handlers.ConfigHandler)

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Wrap the ServeMux with CORS middleware
	handler := c.Handler(mux)
	// serve
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port
	log.Printf("log server listening on %s\n", address)
	log.Fatal(http.ListenAndServe(address, handler))
}
