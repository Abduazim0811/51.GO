package db

import (
	"User/server/models"
	"database/sql"
	"errors"
)

type UserService struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (us *UserService) CreateUser(user *models.User, reply *models.User) error {
	err := us.DB.QueryRow(
		"INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
		user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}
	*reply = *user
	return nil
}

func (us *UserService) GetUser(id int, reply *models.User) error {
	return us.DB.QueryRow("SELECT id, name, email FROM users WHERE id=$1", id).
		Scan(&reply.ID, &reply.Name, &reply.Email)
}

func (us *UserService) UpdateUser(user *models.User, reply *models.User) error {
	result, err := us.DB.Exec(
		"UPDATE users SET name=$1, email=$2 WHERE id=$3",
		user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	*reply = *user
	return nil
}

func (us *UserService) DeleteUser(id int, reply *bool) error {
	result, err := us.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	*reply = true
	return nil
}

func (us *UserService) ListUsers(empty struct{}, reply *[]models.User) error {
	rows, err := us.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return err
		}
		users = append(users, user)
	}
	*reply = users
	return nil
}
