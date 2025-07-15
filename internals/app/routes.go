package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"templ_workout/handlers"
	"templ_workout/internals/auth"
	vAuth "templ_workout/views/auth"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var users = map[string]auth.Login{}

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	fooHandler := &handlers.Foo{
		DB: a.DB,
	}
	authHandler := &handlers.AuthHandler{
		DB: a.DB,
	}
	docHandler := &handlers.Doc{}
	errorPagesHandler := &handlers.ErrorPagesHandler{}

	router.Get("/boo", handlers.Make(fooHandler.HandleFoo))
	router.Get("/", handlers.Make(fooHandler.HandleMoo))
	router.Get("/docs", handlers.Make(docHandler.HandleDocs))
	router.Post("/addUser", fooHandler.HandleAddUser)
	router.Delete("/delete/{email}", fooHandler.HandleDeleteUser)

	router.Get("/register", handlers.Make(authHandler.Register))
	router.Route("/login", func(r chi.Router) {
		r.Use(a.SetUserContext)
		r.Get("/", handlers.Make(authHandler.Login))
	})

	router.Get("/*", handlers.Make(errorPagesHandler.NotFound))

	a.loadAuthRoutes(router, authHandler)

	path, _ := os.Getwd()
	fmt.Println(path)

	fileServer := http.FileServer(http.Dir(path + "/cmd/app/public"))
	router.Handle("/assets/*", fileServer)

	listenAddr := os.Getenv("SERVER_PORT")
	slog.Info("HTTP server started", "server port", listenAddr)
	http.ListenAndServe(listenAddr, router)
	a.router = router
}

func (a *App) loadAuthRoutes(router chi.Router, handler *handlers.AuthHandler) {

	router.Post("/register", handlers.Make(handler.HandleAddUser))

	router.Post("/login", handlers.Make(handler.LoginPost))

	router.Route("/logout", func(r chi.Router) {
		// r.Use(a.Authorize)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				MaxAge:   -1,
				HttpOnly: true,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    "",
				MaxAge:   -1,
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: false,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		})
	})

	router.Route("/protected", func(r chi.Router) {
		r.Use(a.Authorize)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserFromContext(r)
			if user == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}

			handlers.Render(w, r, vAuth.ProtectedPage(user.Name))

		})
	})

}
