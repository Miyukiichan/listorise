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
		Sortable: true,
	})
	columnDTOs = append(columnDTOs, dto.ColumnDTO { 
		Header: "NoteId", 
		Name: "-2", 
		Sortable: true,
	})
	columnDTOs = append(columnDTOs, dto.ColumnDTO { 
		Header: "AssociatedListId", 
		Name: "-3", 
		Sortable: true,
	})
	columnDTOs = append(columnDTOs, dto.ColumnDTO { 
		Header: "ItemId", 
		Name: "-4", 
		Sortable: true,
	}) 
	rows, err := DB().Query("select * from ListColumns where ListId = ?", id)
	if (err != nil && err != sql.ErrNoRows) {
		log.Fatal()
	}
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
			Sortable: true,
		}
		columnDTOs = append(columnDTOs, columnDTO)
	}
	s, err := json.Marshal(columnDTOs)
	if (err != nil) {
		log.Fatal(err)
	}
	listDTO.Columns = template.JS(s)

	itemMap := map[int]map[string]string{}
	
	var listItems []struct {
		Id int
		NoteId sql.NullInt64
		AssociatedListId sql.NullInt64
		NoteName sql.NullString
		AssociatedListName sql.NullString
	}
	rows, err = DB().Query(`
		select i.Id, NoteId, AssociatedListId, n.Name as NoteName, l.Name as AssociatedListName 
		from ListItems as i 
		left join Notes as n on n.Id = i.NoteId 
		left join Lists as l on l.Id = i.AssociatedListId 
		where ListId = ?
	`, id)
	if (err != nil && err != sql.ErrNoRows) {
		log.Fatal()
	}
	err = scan.Rows(&listItems, rows)
	if (err != nil) {
		log.Fatal(err)
	}
	for _, item := range(listItems) {
		itemMap[item.Id] = map[string]string{}
		itemMap[item.Id]["-4"] = strconv.Itoa(item.Id)
		if (item.NoteId != sql.NullInt64{}) {
			itemMap[item.Id]["-1"] = item.NoteName.String
			itemMap[item.Id]["-2"] = strconv.Itoa(int(item.NoteId.Int64))
		} else {
			itemMap[item.Id]["-1"] = item.AssociatedListName.String
			itemMap[item.Id]["-3"] = strconv.Itoa(int(item.AssociatedListId.Int64))
		}
	}

	rows, err = DB().Query(`
		select i.Id as ItemId, c.Id as ColumnId, v.Value 
		from ListValues as v 
		join ListItems as i on v.ListItemId = i.Id 
		join ListColumns as c on i.ListId = c.ListId 
		where c.ListId = ?
	`, id)

	if (err != nil && err != sql.ErrNoRows) {
		log.Fatal()
	}
	for rows.Next() {
		var itemId int
		var columnId int
		var value string
		err = rows.Scan(&itemId, &columnId, &value)
		if err != nil {
			log.Fatal(err)
		}
		itemMap[itemId][strconv.Itoa(columnId)] = value
	}
	var itemList []map[string]string
	for _, value := range itemMap {
		itemList = append(itemList, value)
	}
	s, err = json.Marshal(itemList)
	if (err != nil) {
		log.Fatal(err)
	}
	listDTO.Items = template.JS(s)

	rows.Close()
	tmpl.Execute(w, listDTO)
}