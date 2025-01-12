
///////////////////////////////////////////////////////////////////////
// main.go
///////////////////////////////////////////////////////////////////////
package main

import (
    "log"
    "net/http"

    "logging_microservice/db"
    "logging_microservice/handlers"
)

func main() {
    log.Println("Starting log service")

    // Initialize the database
    if err := db.InitDB("./logging_microservice.db"); err != nil {
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
    http.HandleFunc("/config", handlers.ConfigHandler)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

