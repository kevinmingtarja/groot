package main

import (
	"fmt"
	"net/http"
)

func (s *server) routes() {
	s.router.HandleFunc("/logs/url", s.handleLogsGetByURL).Methods("POST")
	s.router.HandleFunc("/logs", s.handleLogsCreate).Methods("POST")
	s.router.HandleFunc("/logs/{id}", s.handleLogsGetByID).Methods("GET")
	s.router.HandleFunc("/chat", s.handleChatSetID).Methods("POST")
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.RequestURI)

		fmt.Fprintf(w, "Hello!")
	})
}
