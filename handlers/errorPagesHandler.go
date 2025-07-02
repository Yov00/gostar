package handlers

import (
	"net/http"
	"templ_workout/views/layouts"
)

type ErrorPagesHandler struct {
}

func (a *ErrorPagesHandler) NotFound(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, layouts.NotFound())
}
