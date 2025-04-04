// /////////////////////////////////////////////////////////////////////
// handlers/logs.go
// /////////////////////////////////////////////////////////////////////
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"logging_microservice/db"
	"logging_microservice/models"
)

// LogsHandler deals with both POST and GET to /logs
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreateLog(w, r)
	case http.MethodGet:
		handleGetLogs(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleCreateLog(w http.ResponseWriter, r *http.Request) {
	// log.Println("Handling POST request for new log entry")

	var entry models.LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		log.Println("Invalid request payload - ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get or create logger
	var logger models.Logger
	logger, err := db.GetLogger(entry.Logger)
	if err != nil {
		log.Printf("Failed to get/create logger %s: %s\n", logger.Name, err)
		http.Error(w, "Failed to get/create logger", http.StatusInternalServerError)
		return
	}
	entry.LoggerID = logger.ID

	// Check level priority
	if models.LevelPriority[entry.Level] < models.LevelPriority[logger.Level] {
		msg := fmt.Sprintf("Log level for %s too low: %s < %s\n", logger.Name, entry.Level, logger.Level)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(msg))
		return
	}

	// Create the log entry
	// log.Printf("Creating log entry for %s:\n\t%s\n", entry.Logger, entry.Message)
	newEntry, err := db.CreateLog(entry)
	if err != nil {
		log.Println("Failed to create log:", err)
		http.Error(w, "Failed to create log", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEntry)
}

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	// Distinguish between listing all logs or a single log by ID
	path := strings.TrimPrefix(r.URL.Path, "/logs")
	if path == "" || path == "/" {
		// /logs -> list logs
		listLogs(w, r)
	} else {
		// /logs/{id} -> single log
		logIDStr := strings.TrimPrefix(path, "/")
		logID, err := strconv.ParseInt(logIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid log ID", http.StatusBadRequest)
			return
		}
		singleLog(w, logID)
	}
}

func singleLog(w http.ResponseWriter, logID int64) {
	logEntry, err := db.GetSingleLog(logID)
	if err != nil {
		http.Error(w, "Log not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

func listLogs(w http.ResponseWriter, r *http.Request) {
	minTimeStr := r.URL.Query().Get("mintime")
	maxTimeStr := r.URL.Query().Get("maxtime")
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	excludeLoggersStr := r.URL.Query()["excludeLoggers"]
	searchStr := r.URL.Query().Get("search")
	var excludeLoggers []int
	if excludeLoggersStr != nil {
		for _, loggerStr := range excludeLoggersStr {
			loggerId, err := strconv.Atoi(loggerStr)
			if err != nil {
				errMsg := fmt.Sprintf("Invalid logger ID: %s", loggerStr)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
			excludeLoggers = append(excludeLoggers, loggerId)
		}
	}

	logs, err := db.GetLogs(minTimeStr, maxTimeStr, searchStr, offsetStr, limitStr, excludeLoggers)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get logs: %s", err)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
