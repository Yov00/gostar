package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"templ_workout/internals/models"
	"templ_workout/views/foo"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
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
	var users []models.User
	for rows.Next() {
		var user models.User
		rows.Scan(&user.Name, &user.Email)

		users = append(users, user)
	}

	return Render(w, r, foo.Moo(users))
}

func (f *Foo) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(err)
	}

	user.Email = fmt.Sprintf("%d%s", rand.Int(), user.Email)
	_, err = f.DB.Exec("INSERT INTO users(id,name,email, password,createdOn,updatedOn) values($id,$name,$email,$password,$createdOn,$updatedOn)", uuid.New(), user.Name, user.Email, "randomstr", time.Now(), time.Now())
	if err != nil {
		fmt.Println(err)
	}
	Render(w, r, foo.UserContainer(user))
}

func (f *Foo) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	fmt.Println(email)
	if email != "" {
		res, err := f.DB.Exec("DELETE FROM users WHERE email=?", email)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	w.WriteHeader(http.StatusOK)
}
