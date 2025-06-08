package handlers

import (
	"net/http"
	"templ_workout/views/docs"
)

type Doc struct {
}

func (d *Doc) HandleDocs(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, docs.Docs())
}
