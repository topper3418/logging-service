///////////////////////////////////////////////////////////////////////
// db/database.go
///////////////////////////////////////////////////////////////////////
package db

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

// DB is a global reference to the SQL connection
var DB *sql.DB

// InitDB opens the SQLite database
func InitDB(dbFile string) error {
    // make sure the file exists
    if _, err := os.Stat(dbFile); os.IsNotExist(err) {
        log.Println("Database file does not exist, creating")
        file, err := os.Create(dbFile)
        if err != nil {
            return err
        }
        file.Close()
    }
    var err error
    DB, err = sql.Open("sqlite3", dbFile)
    if err != nil {
        return err
    }
    return nil
}

// CreateTables initializes tables if they don't exist
func CreateTables() error {
    log.Println("Ensuring tables exist")

    _, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS logger (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL,
            level TEXT NOT NULL
        );
    `)
    if err != nil {
        return err
    }

    _, err = DB.Exec(`
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
        return err
    }

    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS metadata (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            log_id INTEGER NOT NULL,
            data TEXT NOT NULL,
            FOREIGN KEY(log_id) REFERENCES log(id) ON DELETE CASCADE
        );
    `)
    if err != nil {
        return err
    }

    return nil
}
