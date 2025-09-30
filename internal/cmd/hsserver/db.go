package hsserver

import (
	"database/sql"
	"log"
)

func openDb(dsn string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("Error opening db: %v\n", err)
	}

	// _, err = db.Exec("CREATE TABLE IF NOT EXISTS commands (id INTEGER PRIMARY KEY, command_text TEXT NOT NULL, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	// if err != nil {
	// 	log.Fatalf("Error creating commands table: %v\n", err)
	// }
	return db, nil
}
