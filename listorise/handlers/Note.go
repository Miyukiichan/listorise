package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Miyukiichan/listorise/model/dto"
	"github.com/Miyukiichan/listorise/model/entities"
	"github.com/blockloop/scan"
	"github.com/gorilla/mux"
)

func GetNoteById(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/note.html")
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	noteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	row, err := DB().Query("select * from Notes where Id = ?", noteID)
	if err == sql.ErrNoRows {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatal(err)
	}
	var note entities.Note
	err = scan.Row(&note, row)
	if err != nil {
		log.Fatal(err)
	}
	noteDTO := dto.NoteDTO{
		Id:   noteID,
		Name: note.Name.String,
		Body: note.Body.String,
	}
	row.Close()
	tmpl.Execute(w, noteDTO)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	var noteDTO dto.NoteDTO
	err := json.NewDecoder(r.Body).Decode(&noteDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if noteDTO.Id == 0 {
		http.Error(w, "No Note Id provided", http.StatusBadRequest)
		return
	}
	_, err = DB().Exec("update Notes set Name = ?, Body = ? where Id = ?", noteDTO.Name, noteDTO.Body, noteDTO.Id)
	if err != nil {
		http.Error(w, "Error updating note", http.StatusInternalServerError)
		return
	}
}
