// /////////////////////////////////////////////////////////////////////
// handlers/config.go
// /////////////////////////////////////////////////////////////////////
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
		putHandler(w, r)
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)

	default:
		// If the method is neither PUT nor GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
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
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	type LoggerUpdateRequest struct {
		ID    int    `json:"id"`
		Level string `json:"level"`
	}
	// Handle PUT request to update the logger level
	var loggerConfig LoggerUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&loggerConfig); err != nil {
		errMsg := "Invalid request payload - " + err.Error()
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Update the logger level in the DB
	err := db.UpdateLoggerLevel(loggerConfig.ID, loggerConfig.Level)
	if err != nil {
		log.Println("Failed to update logger level:", err)
		http.Error(w, "Failed to update logger level", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logger level updated successfully"))
	log.Println("Logger level updated successfully")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Handle POST request to create a new logger
	var newLogger models.Logger
	if err := json.NewDecoder(r.Body).Decode(&newLogger); err != nil {
		http.Error(w, "Invalid request payload - ", http.StatusBadRequest)
		return
	}

	// Create the logger in the DB
	newLogger, err := db.GetLogger(newLogger.Name)
	if err != nil {
		log.Println("Failed to create logger:", err)
		http.Error(w, "Failed to create logger", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Logger created successfully"))
}
