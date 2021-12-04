package main

import (
	"fmt"
	"net/http"
)

func (s *server) routes() {
	s.router.HandleFunc("/logs/url", s.getLogsByURLHandler).Methods("POST")
	s.router.HandleFunc("/logs", s.LogHandler).Methods("POST")
	s.router.HandleFunc("/logs/{id}", s.getLogHandler).Methods("GET")
	s.router.HandleFunc("/chat", s.setChatIDHandler).Methods("POST")
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.RequestURI)

		fmt.Fprintf(w, "Hello!")
	})
}
