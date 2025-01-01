package model

import (
	"fmt"
	"log/slog"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
	LoggedIn bool
}

func RegisterUser(registerData RegisterData) error {
	insertQuery := `INSERT INTO user (username, password) VALUES (?, ?)`
	_, err := db.Exec(insertQuery, registerData.Username, registerData.Password)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func GetUserByName(username string) (*User, error) {
	rows, err := db.Query("SELECT id,username, password FROM user where username = ($1)", username)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var user User
	if !rows.Next() {
		err := fmt.Errorf("User not found.")
		return nil, err
	}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id int) (*User, error) {
	rows, err := db.Query("SELECT id,username, password FROM user WHERE id = ($1)", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user User
	if !rows.Next() {
		err := fmt.Errorf("User not found.")
		return nil, err
	}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, username, password FROM user")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
