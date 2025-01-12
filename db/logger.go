package db

import (
    "database/sql"

    "logging_microservice/models"
)

// GetOrCreateLogger retrieves the logger by name or creates a new one
func GetLogger(loggerName string) (models.Logger, error) {
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
func UpdateLoggerLevel(loggerName, newLevel string) error {
    _, err := DB.Exec(`UPDATE logger SET level = ? WHERE name = ?`, newLevel, loggerName)
    return err
}

