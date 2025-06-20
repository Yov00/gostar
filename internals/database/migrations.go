package database

func Migrate() {
	sessionsTableSQL := `
	CREATE TABLE sessions (
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
}
