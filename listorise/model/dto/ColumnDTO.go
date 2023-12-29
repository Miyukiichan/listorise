package dto

// type ColumnEditorListItemDTO struct {
// 	Text string
// 	Value int
// }

// type ColumnEditorOptionsDTO struct {
// 	Format string `json:"format"`
// 	TimePicker bool `json:"timePicker"`
// 	ListItems *[]ColumnEditorListItemDTO `json:"listItems"`
// }

// type ColumnEditorDTO struct {
// 	Type string `json:"name"`
// 	Options *ColumnEditorOptionsDTO `json:"options"`
// }

type ColumnDTO struct {
	Title        string `json:"title"`
	Field        string `json:"field"`
	Sorter       string `json:"sorter"`
	Visible      bool   `json:"visible"`
	HeaderFilter bool   `json:"headerFilter"`
}
