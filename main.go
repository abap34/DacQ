package main

import (
	"github.com/dacq/handler"
	"github.com/dacq/model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("Warn: cannot load env file")
	}

	err := model.InitDB()
	if err != nil {
		log.Fatalf("cannot init database: %v", err)
	}

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
