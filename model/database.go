package model

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
)

var db *sql.DB

func Connect() {
	// Connect to PostgreSQL
	conn, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	db = conn

	// load 'init.sql' and execute it
	file, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(file))
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
