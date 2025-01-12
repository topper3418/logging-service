package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Logger represents an entry in the 'loggers' table
type Logger struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

// LogEntry represents an entry in the 'logs' table
type LogEntry struct {
	ID        int64       `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Logger    string      `json:"logger"`
	LoggerID  int         `json:"logger_id"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Meta      *interface{} `json:"meta,omitempty"`
}

var db *sql.DB
var levelPriority = map[string]int{
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
}

func main() {
	log.Println("Starting log service")
	var err error
	db, err = sql.Open("sqlite3", "./logging_microservice.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	createTables()

	http.HandleFunc("/logs", logsHandler)
    http.HandleFunc("/logs/", logsHandler)
	http.HandleFunc("/config", configHandler)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// createTables initializes tables
func createTables() {
    log.Println("Ensuring Tables")
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS logger (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL,
            level TEXT NOT NULL
        );
    `)
    if err != nil {
        log.Fatal("Failed to create logger table:", err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS log (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp DATETIME NOT NULL,
            logger_id INTEGER NOT NULL,
            level TEXT NOT NULL,
            message TEXT NOT NULL,
            FOREIGN KEY(logger_id) REFERENCES logger(id)
        );
    `)
    if err != nil {
        log.Fatal("Failed to create log table:", err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS metadata (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            log_id INTEGER NOT NULL,
            data TEXT NOT NULL,
            FOREIGN KEY(log_id) REFERENCES log(id) ON DELETE CASCADE
        );
    `)
    if err != nil {
        log.Fatal("Failed to create metadata table:", err)
    }
}

// logsHandler handles creating, retrieving, and querying logs
func logsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /logs")
	if r.Method == http.MethodPost {
		log.Println("Handling POST request to create a log")
		createLog(w, r)
		return
	}

	// Handle GET requests
	if r.Method == http.MethodGet {
		path := strings.TrimPrefix(r.URL.Path, "/logs")
		if path == "" {
			log.Println("Handling GET request to list logs")
			getLogs(w, r) // List logs
		} else {
			logIDStr := strings.TrimPrefix(path, "/")
			logID, err := strconv.ParseInt(logIDStr, 10, 64)
			if err != nil {
				log.Println("Invalid log ID")
				http.Error(w, "Invalid log ID", http.StatusBadRequest)
				return
			}
			log.Printf("Handling GET request for log ID: %d\n", logID)
			getSingleLog(w, logID) // Get single log by ID
		}
		return
	}

	log.Println("Method not allowed")
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// configHandler handles setting logger levels
func configHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /config")
	if r.Method != http.MethodPost {
		log.Println("Method not allowed on /config")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loggerConfig Logger
	err := json.NewDecoder(r.Body).Decode(&loggerConfig)
	if err != nil {
		log.Println("Invalid request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		UPDATE logger SET level = ? WHERE name = ?
	`, loggerConfig.Level, loggerConfig.Name)
	if err != nil {
		log.Println("Failed to update logger level")
		http.Error(w, "Failed to update logger level", http.StatusInternalServerError)
		return
	}

	log.Printf("Updated logger level: %s -> %s\n", loggerConfig.Name, loggerConfig.Level)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logger level updated successfully"))
}

// createLog handles creating a new log and possibly new logger
func createLog(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing create log request")
	var entry LogEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		log.Println("Invalid request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch logger ID
	row := db.QueryRow(`SELECT id, level FROM logger WHERE name = ?`, entry.Logger)
	var loggerID int
	var loggerLevel string
	err = row.Scan(&loggerID, &loggerLevel)
	if err == sql.ErrNoRows {
		loggerLevel = "info" // Default to info if logger doesn't exist
		result, err := db.Exec(`INSERT INTO logger (name, level) VALUES (?, ?)`, entry.Logger, loggerLevel)
		if err != nil {
			log.Println("Failed to create logger")
			http.Error(w, "Failed to create logger", http.StatusInternalServerError)
			return
		}
		lastInsertID, _ := result.LastInsertId()
		loggerID = int(lastInsertID) // Convert int64 to int explicitly
		log.Printf("Created new logger with ID: %d\n", loggerID)
	} else if err != nil {
		log.Println("Failed to fetch logger")
		http.Error(w, "Failed to fetch logger", http.StatusInternalServerError)
		return
	}

	if levelPriority[entry.Level] < levelPriority[loggerLevel] {
		msg := fmt.Sprintf("Log level too low: %s < %s\n", entry.Level, loggerLevel)
		log.Print(msg)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(msg))
		return
	}

	// Insert the log
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	res, err := db.Exec(`
		INSERT INTO log (timestamp, logger_id, level, message)
		VALUES (?, ?, ?, ?)
	`, entry.Timestamp, loggerID, entry.Level, entry.Message)
	if err != nil {
		log.Println("Failed to insert log")
		http.Error(w, "Failed to insert log", http.StatusInternalServerError)
		return
	}

	logID, _ := res.LastInsertId()
	log.Printf("Inserted log with ID: %d\n", logID)

	// Insert metadata if present
	if entry.Meta != nil {
		metaBytes, err := json.Marshal(entry.Meta)
		if err != nil {
			log.Println("Failed to encode metadata")
			http.Error(w, "Failed to encode metadata", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`
			INSERT INTO metadata (log_id, data)
			VALUES (?, ?)
		`, logID, string(metaBytes))
		if err != nil {
			log.Println("Failed to insert metadata")
			http.Error(w, "Failed to insert metadata", http.StatusInternalServerError)
			return
		}
		log.Printf("Inserted metadata for log ID: %d\n", logID)
	}

	entry.ID = logID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// getSingleLog retrieves a single log by ID
func getSingleLog(w http.ResponseWriter, logID int64) {
	log.Printf("Fetching log with ID: %d\n", logID)
	row := db.QueryRow(`
		SELECT log.id, log.timestamp, log.logger_id, logger.name, log.level, log.message
		FROM log INNER JOIN logger ON log.logger_id = logger.id
		WHERE log.id = ?
	`, logID)

	var l LogEntry
	var loggerName string
	err := row.Scan(&l.ID, &l.Timestamp, &l.LoggerID, &loggerName, &l.Level, &l.Message)
	if err != nil {
		log.Println("Log not found")
		http.Error(w, "Log not found", http.StatusNotFound)
		return
	}

	// Fetch metadata
	var metaData string
	metaRow := db.QueryRow(`SELECT data FROM metadata WHERE log_id = ?`, l.ID)
	err = metaRow.Scan(&metaData)
	if err == nil {
		var meta interface{}
		json.Unmarshal([]byte(metaData), &meta)
		l.Meta = &meta
		log.Printf("Fetched metadata for log ID: %d\n", l.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(l)
}

// getLogs handles querying logs with various query params
func getLogs(w http.ResponseWriter, r *http.Request) {
	log.Println("Querying logs with filters")
	minTimeStr := r.URL.Query().Get("mintime")
	maxTimeStr := r.URL.Query().Get("maxtime")
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	includeLoggers := r.URL.Query()["includeLoggers"] // may appear multiple times
	excludeLoggers := r.URL.Query()["excludeLoggers"] // may appear multiple times
	searchStr := r.URL.Query().Get("search")
	queryBuilder := `SELECT log.id, log.timestamp, log.logger_id AS loggerID, logger.name AS logger, log.level, log.message FROM log LEFT JOIN logger ON log.logger_id = logger.id WHERE 1=1`
	args := []interface{}{}

	// Time range filters
	if minTimeStr != "" {
		log.Printf("Applying minTime filter: %s\n", minTimeStr)
		queryBuilder += ` AND log.timestamp >= ?`
		args = append(args, minTimeStr)
	}
	if maxTimeStr != "" {
		log.Printf("Applying maxTime filter: %s\n", maxTimeStr)
		queryBuilder += ` AND log.timestamp <= ?`
		args = append(args, maxTimeStr)
	}

	// Include loggers
	if len(includeLoggers) > 0 {
		log.Printf("Applying includeLoggers filter: %v\n", includeLoggers)
		placeholders := strings.Repeat("?,", len(includeLoggers))
		placeholders = placeholders[:len(placeholders)-1]
		queryBuilder += " AND logger.name IN (" + placeholders + ")"
		for _, l := range includeLoggers {
			args = append(args, l)
		}
	}

	// Exclude loggers
	if len(excludeLoggers) > 0 {
		log.Printf("Applying excludeLoggers filter: %v\n", excludeLoggers)
		placeholders := strings.Repeat("?,", len(excludeLoggers))
		placeholders = placeholders[:len(placeholders)-1]
		queryBuilder += " AND logger.name NOT IN (" + placeholders + ")"
		for _, l := range excludeLoggers {
			args = append(args, l)
		}
	}

	// Search in message
	if searchStr != "" {
		log.Printf("Applying search filter: %s\n", searchStr)
		queryBuilder += ` AND log.message LIKE ?`
		args = append(args, "%"+searchStr+"%")
	}

	// Order and limit
	queryBuilder += " ORDER BY log.timestamp DESC"

	if offsetStr != "" {
		offsetVal, _ := strconv.Atoi(offsetStr)
		log.Printf("Applying offset: %d\n", offsetVal)
		queryBuilder += " LIMIT -1 OFFSET ?"
		args = append(args, offsetVal)
	}
	if limitStr != "" {
		limitVal, _ := strconv.Atoi(limitStr)
		log.Printf("Applying limit: %d\n", limitVal)
		queryBuilder += " LIMIT ?"
		args = append(args, limitVal)
	}

	rows, err := db.Query(queryBuilder, args...)
	if err != nil {
		log.Println("Failed to query logs")
		http.Error(w, "Failed to query logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []LogEntry
	for rows.Next() {
		var l LogEntry
		var loggerName string
		err := rows.Scan(&l.ID, &l.Timestamp, &l.LoggerID, &loggerName, &l.Level, &l.Message)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}
		results = append(results, l)
	}

	log.Printf("Retrieved %d logs\n", len(results))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
