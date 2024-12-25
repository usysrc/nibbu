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
		slog.Error("init.sql: " + err.Error())
		panic("SQL error")
	}
	_, err = db.Exec(string(file))
	if err != nil {
		slog.Error("init.sql:" + err.Error())
		panic("SQL error")
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
