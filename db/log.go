package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"logging_microservice/models"
)

// CreateLog inserts a new log and optionally metadata into the DB
func CreateLog(entry models.LogEntry) (models.LogEntry, error) {
	// If no timestamp was provided, use "now"
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Insert the log
	res, err := DB.Exec(`
        INSERT INTO log (timestamp, logger_id, level, message)
        VALUES (?, ?, ?, ?)
    `, entry.Timestamp, entry.LoggerID, entry.Level, entry.Message)
	if err != nil {
		return models.LogEntry{}, err
	}

	logID, _ := res.LastInsertId()
	entry.ID = logID

	// Insert metadata if present
	if entry.Meta != nil {
		metaBytes, marshalErr := json.Marshal(entry.Meta)
		if marshalErr != nil {
			return entry, marshalErr
		}
		_, metaErr := DB.Exec(`
            INSERT INTO metadata (log_id, data)
            VALUES (?, ?)
        `, logID, string(metaBytes))
		if metaErr != nil {
			return entry, metaErr
		}
	}
	// Now, check database file size and trim if necessary
	if err := trimLogsIfNeeded(); err != nil {
		// Log the error, but don't return it
		log.Printf("error trimming logs: %s", err)
		// log the error
	}
	return entry, nil
}

// trims database when it exceeds 100mb
func trimLogsIfNeeded() error {
	// Get database file name
	rows, err := DB.Query("PRAGMA database_list;")
	if err != nil {
		return err
	}
	defer rows.Close()

	var seq int
	var name string
	var file string
	for rows.Next() {
		err = rows.Scan(&seq, &name, &file)
		if err != nil {
			return err
		}
		if name == "main" {
			break
		}
	}

	// Get file size
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	size := fi.Size()

	if size > 100*1024*1024 {
		// Delete 100 oldest logs
		_, err := DB.Exec(`
			DELETE FROM log WHERE id IN (
			SELECT id FROM log ORDER BY timestamp ASC LIMIT 100
		);`)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetSingleLog retrieves a single log (and metadata) by ID
func GetSingleLog(logID int64) (models.LogEntry, error) {
	row := DB.QueryRow(`
        SELECT log.id, log.timestamp, log.logger_id, logger.name, log.level, log.message
        FROM log
        INNER JOIN logger ON log.logger_id = logger.id
        WHERE log.id = ?
    `, logID)

	var l models.LogEntry
	var loggerName string

	err := row.Scan(&l.ID, &l.Timestamp, &l.LoggerID, &loggerName, &l.Level, &l.Message)
	if err != nil {
		return models.LogEntry{}, err
	}
	l.Logger = loggerName

	// Fetch metadata if present
	var metaData string
	metaRow := DB.QueryRow(`SELECT data FROM metadata WHERE log_id = ?`, l.ID)
	if err := metaRow.Scan(&metaData); err == nil {
		var meta interface{}
		_ = json.Unmarshal([]byte(metaData), &meta)
		l.Meta = &meta
	}

	return l, nil
}

// GetLogs retrieves multiple logs, given various filter parameters
func GetLogs(
	minTimeStr, maxTimeStr, searchStr, offsetStr, limitStr string,
	excludeLoggers []int,
) ([]models.LogEntry, error) {
	queryBuilder := `
        SELECT log.id, log.timestamp, log.logger_id, logger.name, log.level, log.message
        FROM log
        LEFT JOIN logger ON log.logger_id = logger.id
        WHERE 1=1
    `
	args := []interface{}{}

	// Time range filters
	if minTimeStr != "" {
		queryBuilder += ` AND log.timestamp >= ?`
		args = append(args, minTimeStr)
	}
	if maxTimeStr != "" {
		queryBuilder += ` AND log.timestamp <= ?`
		args = append(args, maxTimeStr)
	}

	// Exclude loggers
	if len(excludeLoggers) > 0 {
		placeholders := strings.Repeat("?,", len(excludeLoggers))
		placeholders = placeholders[:len(placeholders)-1]
		queryBuilder += " AND logger.id NOT IN (" + placeholders + ")"
		for _, l := range excludeLoggers {
			args = append(args, l)
		}
	}

	// Search in message
	if searchStr != "" {
		log.Printf("searchStr: %s\n", searchStr)
		log.Printf("searchStr length: %d\n", len(searchStr))
		queryBuilder += ` AND log.message LIKE ?`
		args = append(args, "%"+searchStr+"%")
	}

	// Order and limit
	queryBuilder += " ORDER BY log.timestamp DESC"

	// Offset and limit
	if limitStr != "" {
		queryBuilder += " LIMIT ?"
		args = append(args, limitStr)
	}
	if offsetStr != "" {
		if limitStr == "" {
			return nil, fmt.Errorf("Cannot do an offset without a limit")
		}
		queryBuilder += " OFFSET ?"
		args = append(args, offsetStr)
	}

	rows, err := DB.Query(queryBuilder, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.LogEntry
	for rows.Next() {
		var l models.LogEntry
		var loggerName string
		if err := rows.Scan(&l.ID, &l.Timestamp, &l.LoggerID, &loggerName, &l.Level, &l.Message); err != nil {
			// skip any row that fails
			continue
		}
		l.Logger = loggerName
		results = append(results, l)
	}

	return results, nil
}
