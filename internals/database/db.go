package database

import (
	"database/sql"

    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
}

func (db *DB) NewSqliteDB(filePath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
