package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"templ_workout/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
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
