package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"templ_workout/handlers"
	"templ_workout/internals/auth"
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
	authHandler := &handlers.AuthHandler{}

	docHandler := &handlers.Doc{}
	router.Get("/boo", handlers.Make(fooHandler.HandleFoo))
	router.Get("/", handlers.Make(fooHandler.HandleMoo))
	router.Get("/docs", handlers.Make(docHandler.HandleDocs))
	router.Post("/addUser", fooHandler.HandleAddUser)
	router.Delete("/delete/{email}", fooHandler.HandleDeleteUser)
	router.Get("/login", handlers.Make(authHandler.Login))

	a.loadAuthRoutes(router)

	path, _ := os.Getwd()
	fmt.Println(path)

	fileServer := http.FileServer(http.Dir(path + "/cmd/app/public"))
	router.Handle("/assets/*", fileServer)

	listenAddr := os.Getenv("SERVER_PORT")
	slog.Info("HTTP server started", "server port", listenAddr)
	http.ListenAndServe(listenAddr, router)
	a.router = router
}

func (a *App) loadAuthRoutes(router chi.Router) {
	// orderHandler := &handler.Order{
	// 	Repo: &order.RedisRepo{
	// 		Client: a.rdb,
	// 	},
	// }

	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) < 8 || len(password) < 8 {
			err := http.StatusNotAcceptable
			http.Error(w, "Invalid username/password", err)
			return
		}

		if _, ok := users[username]; ok {
			err := http.StatusConflict
			http.Error(w, "User already exists", err)
			return
		}

		hashedPassword, _ := auth.HashPassword(password)
		users[username] = auth.Login{
			HashedPassword: hashedPassword,
		}

		fmt.Fprintf(w, "User registered successfully!")

	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, ok := users[username]
		if !ok || !auth.CheckPasswordHash(password, user.HashedPassword) {
			err := http.StatusUnauthorized
			http.Error(w, "invalid username/password", err)
			return
		}

		sessionToken := auth.GenerateToken(32)
		csrfToken := auth.GenerateToken(32)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true, //Da ne se pipa ot JS
			SameSite: http.SameSiteStrictMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false, // Da se vzeme ot  js i da se prashta na vseki request
			SameSite: http.SameSiteStrictMode,
		})

		user.SessionToken = sessionToken
		user.CSRFToken = csrfToken
		users[username] = user

		fmt.Fprintf(w, "Login successfully!")

	})

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
		r.Use(Authorize)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			usrname := r.FormValue("username")
			fmt.Fprintf(w, "Здравей %s, Bravo!", usrname)

		})
	})

}

// update and MOVE TO middleware...
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		user, ok := users[username]
		if !ok {
			http.Error(w, "Unauthorized - invalid user", http.StatusUnauthorized)
			return
		}

		st, err := r.Cookie("session_token")
		if err != nil || st.Value == "" || st.Value != user.SessionToken {
			http.Error(w, "Unauthorized - invalid session", http.StatusUnauthorized)
			return
		}

		csrf := r.Header.Get("X-CSRF-TOKEN")
		fmt.Println("----")
		fmt.Println(csrf)
		fmt.Println(user.CSRFToken)
		fmt.Println("----")

		if csrf != user.CSRFToken || csrf == "" {
			http.Error(w, "Unauthorized - invalid CSRF token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
