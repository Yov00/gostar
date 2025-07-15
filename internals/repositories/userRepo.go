package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"templ_workout/internals/models"
	"time"

	"github.com/google/uuid"
)

type UserRepo struct {
	DB *sql.DB
}

func (u *UserRepo) Insert(user models.User) error {
	_, err := u.DB.Exec("INSERT INTO users(id,name,email,password,createdOn,updatedOn) values($id,$name,$email,$password,$createdOn,$updatedOn)", user.Id, user.Name, user.Email, user.Password, user.CreatedOn, user.UpdatedOn)
	fmt.Println(err)
	return err
}

func (u *UserRepo) SelectByEmail(email string) (*models.User, error) {
	var err error
	var user models.User
	var createdOnStr, updatedOnStr string

	row := u.DB.QueryRow("SELECT id, email, name, createdOn, updatedOn, password  from users WHERE email = ?", strings.TrimSpace(email))
	err = row.Scan(&user.Id, &user.Email, &user.Name, &createdOnStr, &updatedOnStr, &user.Password)

	layout := "2006-01-02 15:04:05.999999999-07:00"

	user.CreatedOn, err = time.Parse(layout, createdOnStr)
	if err != nil {
		return nil, fmt.Errorf("invalid createdOn format: %v", err)
	}

	user.UpdatedOn, err = time.Parse(layout, updatedOnStr)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("invalid updatedOn format: %v", err)
	}
	if err == sql.ErrNoRows {
		fmt.Printf("User repo SelectByEmail error: %s", err)
		return nil, nil
	}

	return &user, nil
}

func (u *UserRepo) SelectById(userId uuid.UUID) (*models.User, error) {
	var err error
	var user models.User
	var createdOnStr, updatedOnStr string

	row := u.DB.QueryRow("SELECT id, email, name, createdOn, updatedOn, password  from users WHERE id = ?", userId)
	err = row.Scan(&user.Id, &user.Email, &user.Name, &createdOnStr, &updatedOnStr, &user.Password)

	layout := "2006-01-02 15:04:05.999999999-07:00"

	user.CreatedOn, err = time.Parse(layout, createdOnStr)
	if err != nil {
		return nil, fmt.Errorf("invalid createdOn format: %v", err)
	}

	user.UpdatedOn, err = time.Parse(layout, updatedOnStr)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("invalid updatedOn format: %v", err)
	}
	if err == sql.ErrNoRows {
		fmt.Printf("User repo SelectByEmail error: %s", err)
		return nil, nil
	}

	return &user, nil
}

func (u *UserRepo) GetUserIdByCSRFAndSessionToken(session_token string, csrf string) (*string, error) {
	var err error
	var userId string

	row := u.DB.QueryRow("SELECT user_id from sessions where session_token = ? and csrf_token = ?", strings.TrimSpace(session_token), strings.TrimSpace(csrf))

	err = row.Scan(userId)
	if err != nil {
		fmt.Printf("Could not find user session %v", err)
		return nil, err
	}

	return &userId, nil
}
