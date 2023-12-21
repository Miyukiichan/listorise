package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Miyukiichan/listorise/model"
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
	rows , err := DB().Query("select Id, Name from Lists where Id = ?", id)
	if (err != nil) {
		http.Error(w, "No such list", http.StatusNotFound)
		return;
	}
	list := model.List {
		Id: id,
	}
	var listItems []model.ListItem
	scan.Rows(&listItems, rows)
	s, err := json.Marshal(listItems)
	if (err != nil) {
		log.Fatal(err)
	}
	list.Items = template.JS(s)
	rows.Close()
	if err == sql.ErrNoRows {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, list)
}