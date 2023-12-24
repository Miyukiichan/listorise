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
	var db = handlers.DB()
	db.Exec(`create table Lookups (
		Id integer primary key autoincrement, 
		Name text
	)`)
	db.Exec(`create table LookupOptions (
		Id integer primary key autoincrement, 
		Name text not null,
		LookupId integer not null,
		foreign key (LookupId) references Lookups (Id)
	)`)
	db.Exec(`create table Notes (
		Id integer primary key autoincrement, 
		Name text,
		Body text
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
	db.Exec(`create table ListItems (
		Id integer primary key autoincrement,
		ListId integer not null,
		NoteId integer, 
		AssociatedListId integer, -- This is the list that we link to, not the associated list that this item belongs to
		foreign key (ListId) references Lists (Id),
		foreign key (NoteId) references Notes (Id),
		foreign key (AssociatedListId) references Lists (Id)
	)`)
	db.Exec(`create table ListValues (
		ListItemId integer not null,
		ListColumnId integer not null, 
		Value string, 
		primary key (ListItemId, ListColumnId),
		foreign key (ListItemId) references ListItems (Id),
		foreign key (ListColumnId) references ListColumns (Id)
	)`)

	router := mux.NewRouter()

	router.HandleFunc("/note/{id:[0-9]+}", handlers.GetNoteById).Methods("GET")
	router.HandleFunc("/list/{id:[0-9]+}", handlers.GetListById).Methods("GET")
	router.HandleFunc("/api/listItems/{id:[0-9]+}", handlers.GetListItems).Methods("GET")
	router.HandleFunc("/api/note", handlers.AddNote).Methods("POST")
	router.HandleFunc("/api/list", handlers.AddList).Methods("POST")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", handlers.Config().Port), nil))
}
