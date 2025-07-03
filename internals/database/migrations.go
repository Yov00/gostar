package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) {
	addUsersTable(db)
	addSessionsTable(db)
}

func addSessionsTable(db *sql.DB) {

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

func addUsersTable(db *sql.DB) {

	sessionsTableSQL := `CREATE TABLE users (
	id TEXT PRIMARY KEY, -- UUID stored as TEXT
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	createdOn TEXT NOT NULL, -- Store timestamps as ISO 8601 strings
	updatedOn TEXT NOT NULL
	);`

	_, err := db.Exec(sessionsTableSQL)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Sessions Table created successfully!")
}
