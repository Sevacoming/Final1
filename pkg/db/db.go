package db

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(file string) error {
	if file == "" {
		return errors.New("empty db file")
	}

	var err error
	DB, err = sql.Open("sqlite", file)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date TEXT NOT NULL,
  title TEXT NOT NULL,
  comment TEXT,
  repeat TEXT
 );`)
	return err
}
