package migration

import (
	"database/sql"
	"log"
)

func createTableMessages(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    content TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to createTableMessages: %v", err)
	}
}
