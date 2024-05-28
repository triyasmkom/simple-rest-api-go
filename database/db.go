package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDatabase() {
	var err error

	DB, err = sql.Open("sqlite", "./message.db")
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}
	createTable()
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    content TEXT NOT NULL
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to Create Table: %v", err)
	}
}
