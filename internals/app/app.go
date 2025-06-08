package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"templ_workout/internals/config"
	"templ_workout/internals/database"
	"time"
)

type App struct {
	router http.Handler
	DB     *sql.DB
	config config.Config
}

func NewApp(config config.Config) *App {

	app := &App{
		config: config,
	}

	database := database.DB{}
	db, err := database.NewSqliteDB(config.ConnectionString)

	if err != nil {
		log.Fatal("Failed to initialise db: %w", err)
	}
	app.DB = db

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	port := a.config.ServerPort

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.router,
	}

	err := a.DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	fmt.Println("Starting server on: http://localhost:3000")
	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
