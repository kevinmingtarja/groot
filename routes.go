package main

import (
	"fmt"
	"net/http"
)

func (s *server) routes() {
	// Handler functions don’t actually handle the requests, they return a function that does.
	// This gives us a closure environment in which our handler can operate.
	// Reference: https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
	s.router.HandleFunc("/chat", s.handleChatSetID()).Methods("POST")
	s.router.HandleFunc("/logs", s.handleLogsCreate()).Methods("POST")
	s.router.HandleFunc("/logs/{id}", s.handleLogsGetByID()).Methods("GET")
	s.router.HandleFunc("/logs/url", s.handleLogsGetByURL()).Methods("POST")

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.RequestURI)

		fmt.Fprintf(w, "Hello!")
	})
}
