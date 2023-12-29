package entities

import "database/sql"

type ColumnType int

const (
	Text ColumnType = iota + 1
	Number
	Select
	MultiSelect
	Date
	DateTime
	CreatedDateTime
	LastEditedDateTime
	Checkbox
)

type ListColumn struct {
	Id       int
	Name     string
	Type     ColumnType
	ListId   int
	LookupId sql.NullInt64
}
