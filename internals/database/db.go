package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
}

func (db *DB) NewSqliteDB(filePath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", filePath)
	fmt.Println(filePath)
	if err != nil {
		return nil, err
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
