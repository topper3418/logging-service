package models

// Logger represents an entry in the 'loggers' table
type Logger struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

var LevelPriority = map[string]int{
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
}
