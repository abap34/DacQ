package main

import (
	"github.com/dacq/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.GetHome).Methods("GET")
	r.HandleFunc("/upload-csv", handler.PostUploadCSV).Methods("POST")

	staticHandler := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	http.Handle("/", r)

	// result
	r.HandleFunc("/result", handler.GetResult).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
