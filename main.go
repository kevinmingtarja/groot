package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type Database struct {
	Conn *sql.DB
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	db, err := Initialize()
	if err != nil {
		log.Fatalln("Error establishing connection to database.")
	}
	// initializeBot()

	r := mux.NewRouter()

	r.HandleFunc("/logs/url", db.getLogsByURL).Methods("POST")
	r.HandleFunc("/logs", db.createLog).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "HELLO!")
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("Port is used.")
	}
}

func (db *Database) getLogsByURL(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var b struct{ URL string }
	err := decoder.Decode(&b)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if b.URL == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	logs, err := db.ErrorLogByURL(b.URL)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	res := LogResponse{logs}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(res)
	fmt.Fprintf(w, string(j))
}

func (db *Database) createLog(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var b ErrorLog
	err := decoder.Decode(&b)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = db.CreateErrorLog(&b)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Call bot to send message

	fmt.Fprintf(w, "Success")
}
