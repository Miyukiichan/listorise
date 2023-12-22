package entities

import "database/sql"

type Note struct {
	ID int
	Name sql.NullString
	Body sql.NullString
}