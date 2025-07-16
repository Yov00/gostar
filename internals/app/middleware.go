package app

import (
	"context"
	"net/http"
	"templ_workout/internals/models"
	"templ_workout/internals/repositories"
)

func (a *App) SetUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		var err error

		st, err := r.Cookie("session_token")
		if err != nil || st.Value == "" {
			user = nil
		}
		ct, err := r.Cookie("csrf_token")
		if err != nil || ct.Value == "" {
			user = nil
		}

		if st != nil && ct != nil {
			sessionRepo := repositories.SessionRepo{DB: a.DB}
			session, err := sessionRepo.GetSessionByCSRFAndSessionToken(st.Value, ct.Value)
			if err != nil || session == nil {
				user = nil
			}
			userRepo := repositories.UserRepo{DB: a.DB}
			userId := session.UserID
			user, err = userRepo.SelectById(userId)
			if err != nil || user == nil {
				user = nil
			}
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *App) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		st, err := r.Cookie("session_token")
		if err != nil || st.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ct, err := r.Cookie("csrf_token")
		if err != nil || ct.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if r.Method != http.MethodGet {
			headerCT := r.Header.Get("X-CSRF-TOKEN")
			if headerCT == "" || headerCT != ct.Value {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		sessionRepo := repositories.SessionRepo{DB: a.DB}
		session, err := sessionRepo.GetSessionByCSRFAndSessionToken(st.Value, ct.Value)
		if err != nil || session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userRepo := repositories.UserRepo{DB: a.DB}
		userId := session.UserID
		user, err := userRepo.SelectById(userId)
		if err != nil || user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
