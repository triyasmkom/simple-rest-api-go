package migration

import (
	"database/sql"
	"log"
)

func createTableAuthUser(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS auth_user (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    username VARCHAR(50) UNIQUE NOT NULL,
	    first_name VARCHAR(200) DEFAULT NULL,
	    last_name VARCHAR(200) DEFAULT NULL,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to createTableAuthUser: %v", err)
	}
}

func alterTriggerTableAuthUser(db *sql.DB) {
	query := `
	CREATE TRIGGER IF NOT EXISTS update_auth_user_timestamp
			AFTER UPDATE ON auth_user
			FOR EACH ROW
			BEGIN
				UPDATE auth_user SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
			END;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to alterTriggerTableAuthUser: %v", err)
	}
}

func alterColumnPasswordAuthUser(db *sql.DB) {
	query := `ALTER TABLE auth_user ADD column password varchar(255);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to alterColumnPasswordAuthUser: %v", err)
	}
}
