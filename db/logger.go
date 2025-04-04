package db

import (
	"database/sql"
	"log"

	"logging_microservice/models"
)

// GetOrCreateLogger retrieves the logger by name or creates a new one
func GetLogger(loggerName string) (models.Logger, error) {
	// log.Printf("searching for logger %s", loggerName)
	row := DB.QueryRow(`SELECT id, level FROM logger WHERE name = ?`, loggerName)

	var loggerID int
	var loggerLevel string

	err := row.Scan(&loggerID, &loggerLevel)
	if err == sql.ErrNoRows {
		// If logger doesn't exist, create one with a default level of "info"
		loggerLevel = "info"
		res, insertErr := DB.Exec(`INSERT INTO logger (name, level) VALUES (?, ?)`, loggerName, loggerLevel)
		if insertErr != nil {
			return models.Logger{}, insertErr
		}
		lastInsertID, _ := res.LastInsertId()
		loggerID = int(lastInsertID)
		log.Printf("Logger %s created successfully", loggerName)
		return models.Logger{
			ID:    loggerID,
			Name:  loggerName,
			Level: loggerLevel,
		}, nil
	}
	if err != nil {
		return models.Logger{}, err
	}

	return models.Logger{
		ID:    loggerID,
		Name:  loggerName,
		Level: loggerLevel,
	}, nil
}

// UpdateLoggerLevel updates the named logger to a new level
func UpdateLoggerLevel(loggerId int, newLevel string) error {
	_, err := DB.Exec(`UPDATE logger SET level = ? WHERE id = ?`, newLevel, loggerId)
	return err
}

// ListLoggers lists all loggers and their levels.
func ListLoggers() ([]models.Logger, error) {
	// log.Printf("Listing loggers")
	rows, err := DB.Query(`SELECT id, name, level FROM logger`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loggers []models.Logger
	for rows.Next() {
		var logger models.Logger
		if err := rows.Scan(&logger.ID, &logger.Name, &logger.Level); err != nil {
			return nil, err
		}
		loggers = append(loggers, logger)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loggers, nil
}
