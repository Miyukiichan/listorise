package entities

import "database/sql"

type ListValue struct {
	ListItemId   int
	ListColumnId int
	Value        sql.NullString
}
