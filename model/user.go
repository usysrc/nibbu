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
	insertQuery := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err := db.Exec(insertQuery, registerData.Username, registerData.Password)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func GetUserByName(username string) (*User, error) {
	rows, err := db.Query("SELECT id,username, password FROM users where username = ($1)", username)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var user User
	if !rows.Next() {
		err := fmt.Errorf("User not found.")
		slog.Error(err.Error())
		return nil, err
	}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id int) (*User, error) {
	rows, err := db.Query("SELECT id,username, password FROM users where id = ($1)", id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var user User
	if !rows.Next() {
		err := fmt.Errorf("User not found.")
		slog.Error(err.Error())
		return nil, err
	}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &user, nil
}
