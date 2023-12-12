package models

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/sqlite"
)

var Db *sql.DB

var schema = `
CREATE TABLE IF NOT EXISTS monitoring (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	seq_id INTEGER NOT NULL,
    source TEXT NOT NULL,
    status TEXT NOT NULL,
    response_time REAL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
`

func InitDB() {
	// Open SQLite database
	var err error
	Db, err = sql.Open("sqlite", "database/database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	_, err = Db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}
