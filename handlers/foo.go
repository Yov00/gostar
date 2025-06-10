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

func (f *Foo) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("it hits hard")
	decoder := json.NewDecoder(r.Body)
	var user models.UserDTO
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(user)

	user.Email = fmt.Sprintf("%d%s", rand.Int(), user.Email)
	_, err = f.DB.Exec("INSERT INTO users(name,email) values($name,$email)", user.Name, user.Email)
	if err != nil {
		fmt.Println(err)
	}

	Render(w, r, foo.UserContainer(user))

}
