package entities

import "database/sql"

type ListItem struct {
	Id int
	ListId int
	NoteId sql.NullInt64
	AssociatedListId sql.NullInt64
}