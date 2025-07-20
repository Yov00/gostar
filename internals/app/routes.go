package app

import (
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
	router.Use(a.SetUserContext)

	fooHandler := &handlers.Foo{
		DB: a.DB,
	}
	authHandler := &handlers.AuthHandler{
		DB: a.DB,
	}
	docHandler := &handlers.Doc{}
	errorPagesHandler := &handlers.ErrorPagesHandler{}

	router.Get("/", handlers.Make(fooHandler.HandleMoo))
	router.Get("/docs", handlers.Make(docHandler.HandleDocs))
	router.Post("/addUser", fooHandler.HandleAddUser)
	router.Delete("/delete/{email}", fooHandler.HandleDeleteUser)

	router.Get("/*", handlers.Make(errorPagesHandler.NotFound))

	a.loadAuthRoutes(router, authHandler)

	path, _ := os.Getwd()

	fileServer := http.FileServer(http.Dir(path + "/cmd/app/public"))
	router.Handle("/assets/*", fileServer)

	listenAddr := os.Getenv("SERVER_PORT")
	slog.Info("HTTP server started", "server port", listenAddr)
	http.ListenAndServe(listenAddr, router)
	a.router = router
}

func (a *App) loadAuthRoutes(router chi.Router, handler *handlers.AuthHandler) {

	router.Route("/register", func(r chi.Router) {

		r.Get("/", handlers.Make(handler.Register))
		r.Post("/", handlers.Make(handler.HandleAddUser))

	})
	router.Route("/login", func(r chi.Router) {
		r.Post("/", handlers.Make(handler.LoginPost))
		r.Get("/", handlers.Make(handler.Login))
	})
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
