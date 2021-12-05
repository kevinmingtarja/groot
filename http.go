package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func respond(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "respond"))
	}
}

func respondErr(w http.ResponseWriter, r *http.Request, err error, code int) {
	errObj := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(errObj)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "respondErr"))
	}
}
