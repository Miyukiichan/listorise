package entities

import "database/sql"

type List struct {
	Id   int
	Name sql.NullString
}
