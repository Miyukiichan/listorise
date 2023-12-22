package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Miyukiichan/listorise/model/dto"
	"github.com/Miyukiichan/listorise/model/entities"
	"github.com/blockloop/scan"
	"github.com/gorilla/mux"
)

func GetListById(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/list.html")
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Get list name + check if list exists
	var listName string
	err = DB().QueryRow("select Name from Lists where Id = ?", id).Scan(&listName)
	if err == sql.ErrNoRows {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatal(err)
	}
	listDTO := dto.ListDTO {
		Id: id,
	}
	listDTO.Name = listName

	// Get Columns + include special columns for the name, associate list and note for navigation
	var columns []entities.ListColumn
	var columnDTOs []dto.ColumnDTO
	columnDTOs = append(columnDTOs, dto.ColumnDTO { 
		Header: "Name", 
		Name: "-1", 
	})
	columnDTOs = append(columnDTOs, dto.ColumnDTO { 
		Header: "NoteId", 
		Name: "-2", 
	})
	columnDTOs = append(columnDTOs, dto.ColumnDTO { 
		Header: "AssociatedListId", 
		Name: "-3", 
	})
	rows, err := DB().Query("select * from ListColumns where ListId = ?", id)
	if err == nil || err != sql.ErrNoRows {
		err = scan.Rows(&columns, rows)
		if (err != nil) {
			log.Fatal(err)
		}
		for _, column := range columns {
			columnDTO := dto.ColumnDTO {
				// Use the column Id as this is a unique name
				// The name is never visible - it's an internal reference for the data objects in toast ui
				Name: strconv.Itoa(column.Id), 
				Header: column.Name,
				Editor: nil,
			}
			columnDTOs = append(columnDTOs, columnDTO)
		}
	}
	s, err := json.Marshal(columnDTOs)
	if (err != nil) {
		log.Fatal(err)
	}
	listDTO.Columns = template.JS(s)
	listDTO.Items = "[]" // Default value to prevent invalid js
	rows.Close()
	tmpl.Execute(w, listDTO)
}