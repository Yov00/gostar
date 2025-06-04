package handlers

import (
	"net/http"
	"templ_workout/views/foo"
)

type Foo struct {
}

func (f *Foo) HandleFoo(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, foo.Index())
}
