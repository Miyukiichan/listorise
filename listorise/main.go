package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Miyukiichan/listorise/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var db = handlers.DB();
	db.Exec(`create table Lookups (
		Id integer primary key autoincrement, 
		Name text
	)`)
	db.Exec(`create table LookupOptions (
		Id integer primary key autoincrement, 
		LookupId integer not null,
		Name text not null,
		foreign key (LookupId) references Lookups (Id)
	)`)
	db.Exec(`create table Notes (
		Id integer primary key autoincrement, 
		Body text, 
		Title text
	)`)
	db.Exec(`create table Lists (
		Id integer primary key autoincrement, 
		Name text
	)`)
	db.Exec(`create table ListColumns (
		Id integer primary key autoincrement, 
		Name text not null, 
		Type integer not null, 
		ListId integer not null, 
		LookupId integer,
		foreign key (ListId) references Lists (Id),
		foreign key (LookupId) references Lookups (Id)
	)`)
	db.Exec(`create table ListValues (
		Id integer primary key autoincrement, 
		Value string, 
		ListColumnId integer not null, 
		NoteId integer, 
		ListId integer,
		foreign key (ListColumnId) references ListColumns (Id),
		foreign key (NoteId) references Notes (Id),
		foreign key (ListId) references Lists (Id) -- This is the list that we link to, not the parent list
	)`)

	router := mux.NewRouter()

	// Route to handle requests for a specific note based on ID
	router.HandleFunc("/note/{id:[0-9]+}", handlers.GetNoteById).Methods("GET")
	router.HandleFunc("/list/{id:[0-9]+}", handlers.GetListById).Methods("GET")

	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", handlers.Config().Port), nil))
}

