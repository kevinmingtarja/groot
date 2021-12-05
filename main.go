package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load(".env")
	if err != nil {
		return errors.Wrap(err, "environment variables")
	}

	db, err := setupDatabase()
	if err != nil {
		return errors.Wrap(err, "setup database")
	}

	bot, err := setupBot()
	if err != nil {
		return errors.Wrap(err, "setup bot")
	}

	srv := newServer()
	srv.db, srv.bot = db, bot

	return http.ListenAndServe(":8080", srv)
}

// The server type contains the dependencies of our server.
type server struct {
	router *mux.Router
	db     *sql.DB
	bot    *tgbotapi.BotAPI
}

// newServer instantiates a server type and sets up its routes.
// Dependencies are not set up here so that it is easier to test.
func newServer() *server {
	srv := &server{
		router: mux.NewRouter(),
	}
	srv.routes()
	return srv
}

// Implementing ServeHTTP turns the server type into a http.Handler.
// Hence, server can be used wherever http.Handler can (e.g. http.ListenAndServe).
// Inside, we simply pass the execution to the router.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleLogsGetByURL() http.HandlerFunc {
	type request struct {
		URL string
	}
	type response struct {
		Logs []ErrorLog `json:"logs"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var req request
		err := decoder.Decode(&req)
		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}
		if req.URL == "" {
			respondErr(w, r, errors.New("URL cannot be empty"), http.StatusBadRequest)
			return
		}

		logs, err := s.ErrorLogByURL(req.URL)
		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}
		res := response{logs}

		//w.Header().Set("Content-Type", "application/json")
		//j, _ := json.Marshal(res)
		//fmt.Fprintf(w, string(j))
		respond(w, r, res, http.StatusOK)
	}
}

func (s *server) handleLogsCreate() http.HandlerFunc {
	type request ErrorLog
	type response struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var req request
		err := decoder.Decode(&req)
		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}
		errorLog := (ErrorLog)(req)

		id, err := s.CreateErrorLog(&errorLog)
		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}
		errorLog.ID = id

		// Call bot to send message
		chatID, err := s.ChatID(errorLog.AppName)
		if err != nil {
			log.Println(err)
			if chatID == 0 {
				respondErr(w, r, errors.New("No Chat ID under the given app name found. "+
					"Please map it first using the /chat endpoint."), http.StatusBadRequest)
				return
			}
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		err = s.SendErrorMessage(chatID, &errorLog)
		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		res := response{"Error successfully logged"}

		respond(w, r, res, http.StatusCreated)
	}
}

func (s *server) handleChatSetID() http.HandlerFunc {
	type request Chat
	type response struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var req request
		err := decoder.Decode(&req)
		if err != nil {
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}
		chat := (Chat)(req)

		err = s.SetChatID(&chat)
		if err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
			return
		}

		res := response{"Chat ID succesfully mapped"}

		respond(w, r, res, http.StatusOK)
	}
}

func (s *server) handleLogsGetByID() http.HandlerFunc {
	type response struct {
		Log ErrorLog `json:"log"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		i, err := strconv.Atoi(id)
		if err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
			return
		}

		errLog, err := s.ErrorLog(i)
		if err != nil {
			log.Println(err)
			if errLog.ID == 0 {
				respondErr(w, r, err, http.StatusBadRequest)
				return
			}
			respondErr(w, r, err, http.StatusInternalServerError)
			return
		}

		res := response{errLog}
		//w.Header().Set("Content-Type", "application/json")
		//j, _ := json.Marshal(res)
		//fmt.Fprintf(w, string(j))
		respond(w, r, res, http.StatusOK)
	}
}
