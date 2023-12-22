package dto

import "html/template"

type ListDTO struct {
	Id int
	Name string
	Items template.JS
	Columns template.JS
}