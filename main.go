
///////////////////////////////////////////////////////////////////////
// main.go
///////////////////////////////////////////////////////////////////////
package main

import (
    "log"
    "os"
    "net/http"

    "github.com/joho/godotenv"

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
    if err := db.InitDB("./data/logs.db"); err != nil {
        log.Fatal("Database initialization failed:", err)
    }
    defer db.DB.Close()

    // Create tables if not present
    if err := db.CreateTables(); err != nil {
        log.Fatal("Failed to create tables:", err)
    }

    // Register routes
    http.HandleFunc("/logs", handlers.LogsHandler)
    http.HandleFunc("/logs/", handlers.LogsHandler)
    http.HandleFunc("/loggers", handlers.ConfigHandler)

    // serve
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    address := ":" + port
    log.Printf("log server listening on %s\n", address)
    log.Fatal(http.ListenAndServe(address, nil))
}

