package handlers

import (
	"database/sql"
	"log"
)

func DB() *sql.DB {
	var _db *sql.DB
	config := Config()
	var err error
	_db, err = sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	return _db
}
