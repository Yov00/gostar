package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) {
	sessionsTableSQL := `
	CREATE TABLE IF NOT EXISTS sessions  (
	id INTEGER PRIMARY KEY AUTOINCREMENT, -- or UUID
	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	session_token TEXT NOT NULL UNIQUE,
	csrf_token TEXT NOT NULL,
	user_agent TEXT,
	ip_address TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	expires_at DATETIME
	);

	`

	_, err := db.Exec(sessionsTableSQL)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Sessions Table created successfully!")

}
