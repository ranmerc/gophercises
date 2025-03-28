package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Open() error {
	var err error

	DB, err = sql.Open("sqlite3", "todos.db")
	if err != nil {
		return err
	}

	sqlStatement := `
		CREATE TABLE IF NOT EXISTS todos (
		id INTEGER NOT NULL PRIMARY KEY,
		time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		description TEXT NOT NULL
	);`

	if _, err := DB.Exec(sqlStatement); err != nil {
		return err
	}

	return nil
}

func Close() {
	DB.Close()
}
