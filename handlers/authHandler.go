package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	intAuth "templ_workout/internals/auth"
	"templ_workout/internals/models"
	"templ_workout/internals/repositories"
	vAuth "templ_workout/views/auth"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	DB *sql.DB
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, vAuth.Login())
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, vAuth.Register())
}
func (a *AuthHandler) HandleAddUser(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	fmt.Println(email)
	fmt.Println(confirmPassword)

	if len(name) < 3 || len(password) < 6 {
		err := http.StatusNotAcceptable
		http.Error(w, "Invalid username/password", err)
		return nil
	}

	userRepo := repositories.UserRepo{DB: a.DB}

	user, err := userRepo.SelectByEmail(email)
	if err != nil {
		http.Error(w, "failed to check for conflict", http.StatusInternalServerError)
		fmt.Println(err)
		return err
	}

	if user.Email != "" {
		err := http.StatusConflict
		http.Error(w, "User already exists", err)
		return nil
	}

	hashedPassword, _ := intAuth.HashPassword(password)
	err = userRepo.Insert(models.User{
		Id:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	})

	if err != nil {
		err := http.StatusInternalServerError
		http.Error(w, "Registration failed, server error!", err)
		return nil
	}

	fmt.Fprintf(w, "User registered successfully!")
	http.Redirect(w, r, "/login", http.StatusOK)
	return nil

}
