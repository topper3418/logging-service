package models

import (
    "time"
)

type LogEntry struct {
	ID        int64       `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Logger    string      `json:"logger"`
	LoggerID  int         `json:"logger_id"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Meta      *interface{} `json:"meta,omitempty"`
}

