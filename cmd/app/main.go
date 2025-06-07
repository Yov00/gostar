package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"templ_workout/handlers"
	"templ_workout/internals/database"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get working dir %w", err)
	}

	dbPath := filepath.Join(basePath, "db", "app.db")
	fmt.Println("\n")
	fmt.Println("\n")
	fmt.Println("\n")

	fmt.Println(dbPath)

	fmt.Println("\n")
	fmt.Println("\n")
	fmt.Println("\n")

	database := database.DB{}
	db, err := database.NewSqliteDB(dbPath)
	if err != nil {
		log.Fatal("Failed to initialise db: %w", err)
	}

	con, _ := db.Conn(context.Background())
	con.ExecContext(context.Background(), `
		CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		age INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		`)

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()
	fooHandler := &handlers.Foo{}
	router.Get("/boo", handlers.Make(fooHandler.HandleFoo))

	router.Get("/sex", handlers.Make(fooHandler.HandleMoo))

	path, _ := os.Getwd()
	fmt.Println(path)

	fileServer := http.FileServer(http.Dir(path + "/cmd/app/public"))
	router.Handle("/assets/*", fileServer)

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	http.ListenAndServe(listenAddr, router)

}
