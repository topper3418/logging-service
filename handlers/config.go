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
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

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
}

