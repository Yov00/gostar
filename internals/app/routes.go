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
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// router.Route("/orders", a.loadOrderRoutes)
	fooHandler := &handlers.Foo{}
	router.Get("/boo", handlers.Make(fooHandler.HandleFoo))

	router.Get("/sex", handlers.Make(fooHandler.HandleMoo))

	path, _ := os.Getwd()
	fmt.Println(path)

	fileServer := http.FileServer(http.Dir(path + "/cmd/app/public"))
	router.Handle("/assets/*", fileServer)

	listenAddr := os.Getenv("SERVER_PORT")
	slog.Info("HTTP server started", "server port", listenAddr)
	http.ListenAndServe(listenAddr, router)

	a.router = router
}
