package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) {
	addUsersTable(db)
	addSessionsTable(db)
	addFilesTable(db)
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

	sessionsTableSQL := `CREATE TABLE IF NOT EXISTS users (
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

	fmt.Println("Users Table created successfully!")
}

func addFilesTable(db *sql.DB) {
	filesTableSQL := `CREATE TABLE IF NOT EXISTS files (
	id TEXT PRIMARY KEY,
	file_name TEXT NOT NULL,
	original_file_name TEXT NOT NULL,
	file_path TEXT UNIQUE NOT NULL,
	created_on TEXT NOT NULL DEFAULT (datetime('now')),
	updated_on TEXT NOT NULL DEFAULT (datetime('now'))
	);`

	_, err := db.Exec(filesTableSQL)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("files Table created successfully!")
}
