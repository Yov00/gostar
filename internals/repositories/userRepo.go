package repositories

import (
	"database/sql"
	"fmt"
	"templ_workout/internals/models"
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

	row := u.DB.QueryRow("SELECT * from users WHERE email = $email", email)

	if err == sql.ErrNoRows {
		fmt.Println(err)
		return nil, nil
	}

	err = row.Scan(user.Id, user.Email, user.Name, user.CreatedOn, user.CreatedOn, user.Password)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil

}
