package migration

import (
	"database/sql"
	"fmt"
	"log"
)

func Migration(db *sql.DB) {
	// Create Table Messages
	createTableMessages(db)

	// Create Table Auth User
	createTableAuthUser(db)

	// Add Trigger column updated at Table Auth User
	alterTriggerTableAuthUser(db)

	// Add column Password
	exists, err := columnExists(db, "auth_user", "password")
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		alterColumnPasswordAuthUser(db)
	}
}

func columnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var cid int
	var name, ctype string
	var dflt_value sql.NullString
	var notnull, pk int

	for rows.Next() {
		err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
		if err != nil {
			return false, err
		}
		if name == columnName {
			return true, nil
		}
	}
	return false, nil
}
