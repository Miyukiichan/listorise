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
	db.Exec("create table Notes (Id integer primary key autoincrement, Body text, Title text)")
	db.Exec("create table ListItems (Id integer, Name text)")

	router := mux.NewRouter()

	// Route to handle requests for a specific note based on ID
	router.HandleFunc("/note/{id:[0-9]+}", handlers.GetNoteById).Methods("GET")
	router.HandleFunc("/list/{id:[0-9]+}", handlers.GetListById).Methods("GET")

	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", handlers.Config().Port), nil))
}

