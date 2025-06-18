package handlers

import (
	"net/http"
	"templ_workout/views/auth"
)

type AuthHandler struct {
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, auth.Login())
}
