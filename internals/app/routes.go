package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"templ_workout/handlers"
	"templ_workout/internals/auth"
	"templ_workout/internals/repositories"
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
	router.Get("/login", handlers.Make(authHandler.Login))
	router.Get("/register", handlers.Make(authHandler.Register))

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
		r.Use(Authorize)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: false,
			})

			username := r.FormValue("username")
			user, _ := users[username]
			user.SessionToken = ""
			user.CSRFToken = ""
			users[username] = user

			fmt.Fprintf(w, "Logout susccessul!")
		})
	})

	router.Route("/protected", func(r chi.Router) {
		r.Use(a.Authorize)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			usrname := r.FormValue("username")
			fmt.Fprintf(w, "Здравей %s, Bravo!", usrname)

		})
	})

}

func (a *App) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		st, err := r.Cookie("session_token")
		if err != nil || st.Value == "" {
			http.Error(w, "Unauthorized - invalid session", http.StatusUnauthorized)
			return
		}
		ct, err := r.Cookie("csrf_token")
		if err != nil || ct.Value == "" {
			http.Error(w, "Unauthorized - invalid csrf token", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodGet {
			headerCT := r.Header.Get("X-CSRF-TOKEN")
			if headerCT == "" || headerCT != ct.Value {
				http.Error(w, "Unauthorized - invalid csrf token", http.StatusUnauthorized)
				return
			}
		}

		userRepo := repositories.UserRepo{DB: a.DB}
		userId, err := userRepo.GetUserIdByCSRFAndSessionToken(st.Value, ct.Value)
		if err != nil || userId == nil {
			http.Error(w, "Unauthorized - Invalid session", http.StatusUnauthorized)
			return
		}

		user, err := userRepo.SelectById(*userId)
		if err != nil || user == nil {
			http.Error(w, "Unauthorized - Invalid session", http.StatusUnauthorized)
			return
		}

		csrf := r.Header.Get("X-CSRF-TOKEN")
		fmt.Println("----")
		fmt.Println(csrf)
		fmt.Println("----")

		next.ServeHTTP(w, r)
	})
}
