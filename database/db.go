package database

import (
	"database/sql"
	"log"
	mg "rest-api-gorilla/database/migration"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDatabase() {
	var err error

	DB, err = sql.Open("sqlite", "./message.db")
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	mg.Migration(DB)
}

func Query(query string, args ...interface{}) ([]interface{}, error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Membuat slice untuk nilai-nilai yang akan di-scan
	values := make([]interface{}, len(columns))
	for i := range values {
		var value interface{}
		values[i] = &value
	}

	var results []interface{}
	for rows.Next() {
		// Scanning nilai ke dalam slice
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		// Membuat map untuk pasangan kolom-nilai
		result := make(map[string]interface{})
		for i, col := range columns {
			result[col] = *(values[i].(*interface{}))
		}

		results = append(results, result)
	}

	return results, nil
}
