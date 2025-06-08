package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"templ_workout/internals/models"
	"templ_workout/views/foo"
)

type Foo struct {
	DB *sql.DB
}

func (f *Foo) HandleFoo(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, foo.Index())
}

func (f *Foo) HandleMoo(w http.ResponseWriter, r *http.Request) error {
	rows, err := f.DB.Query("select name,email from users")
	if err != nil {
		log.Fatal("failed to query users")
	}
	defer rows.Close()
	var users []models.UserDTO
	for rows.Next() {
		var user models.UserDTO
		rows.Scan(&user.Name, &user.Email)

		users = append(users, user)
	}

	return Render(w, r, foo.Moo(users))
}
