package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"templ_workout/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	fooHandler := &handlers.Foo{
		DB: a.DB,
	}
	authHandler := &handlers.AuthHandler{}

	docHandler := &handlers.Doc{}
	router.Get("/boo", handlers.Make(fooHandler.HandleFoo))
	router.Get("/", handlers.Make(fooHandler.HandleMoo))
	router.Get("/docs", handlers.Make(docHandler.HandleDocs))
	router.Post("/addUser", fooHandler.HandleAddUser)
	router.Delete("/delete/{email}", fooHandler.HandleDeleteUser)
	router.Get("/login", handlers.Make(authHandler.Login))

	path, _ := os.Getwd()
	fmt.Println(path)

	fileServer := http.FileServer(http.Dir(path + "/cmd/app/public"))
	router.Handle("/assets/*", fileServer)

	listenAddr := os.Getenv("SERVER_PORT")
	slog.Info("HTTP server started", "server port", listenAddr)
	http.ListenAndServe(listenAddr, router)

	a.router = router
}
