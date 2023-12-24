package handlers

import (
	"database/sql"
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
		Id:   note.ID,
		Name: note.Name.String,
		Body: note.Body.String,
	}
	row.Close()
	tmpl.Execute(w, noteDTO)
}
