///////////////////////////////////////////////////////////////////////
// handlers/config.go
///////////////////////////////////////////////////////////////////////
package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "logging_microservice/db"
    "logging_microservice/models"
)

// ConfigHandler handles setting logger levels
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPut:
        // Handle PUT request to update the logger level
        var loggerConfig models.Logger
        if err := json.NewDecoder(r.Body).Decode(&loggerConfig); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Update the logger level in the DB
        err := db.UpdateLoggerLevel(loggerConfig.Name, loggerConfig.Level)
        if err != nil {
            log.Println("Failed to update logger level:", err)
            http.Error(w, "Failed to update logger level", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Logger level updated successfully"))

    case http.MethodGet:
        // Handle GET request to list all loggers
        loggers, err := db.ListLoggers()
        if err != nil {
            log.Println("Failed to retrieve loggers:", err)
            http.Error(w, "Failed to retrieve loggers", http.StatusInternalServerError)
            return
        }

        // Return the loggers as JSON
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(loggers); err != nil {
            log.Println("Failed to encode loggers:", err)
            http.Error(w, "Failed to encode loggers", http.StatusInternalServerError)
            return
        }

    default:
        // If the method is neither PUT nor GET
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

