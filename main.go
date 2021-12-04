package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Env struct {
	DB  *sql.DB
	Bot *tgbotapi.BotAPI
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	db, err := InitializeDB()
	if err != nil {
		return err
	}
	bot, err := InitializeBot()
	if err != nil {
		return err
	}
	env := &Env{db, bot}

	r := mux.NewRouter()

	r.HandleFunc("/logs/url", env.getLogsByURLHandler).Methods("POST")
	r.HandleFunc("/logs", env.LogHandler).Methods("POST")
	r.HandleFunc("/logs/{id}", env.getLogHandler).Methods("GET")
	r.HandleFunc("/chat", env.setChatIDHandler).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.RequestURI)

		fmt.Fprintf(w, "Hello!")
	})

	return http.ListenAndServe(":8080", r)
}

func (env *Env) getLogsByURLHandler(w http.ResponseWriter, r *http.Request) {
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

	logs, err := env.ErrorLogByURL(b.URL)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	res := LogsResponse{logs}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(res)
	fmt.Fprintf(w, string(j))
}

func (env *Env) LogHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var errorLog ErrorLog
	err := decoder.Decode(&errorLog)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	id, err := env.CreateErrorLog(&errorLog)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	errorLog.ID = id

	// Call bot to send message
	chatID, err := env.ChatID(errorLog.AppName)
	if err != nil {
		log.Println(err)
		if chatID == 0 {
			http.Error(w, "No Chat ID under the given app name found. Please map it first using the /chat endpoint.", 400)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = env.SendErrorMessage(chatID, &errorLog)

	fmt.Fprintf(w, "Success")
}

func (env *Env) setChatIDHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var b Chat
	err := decoder.Decode(&b)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = env.SetChatID(&b)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprintf(w, "Chat ID succesfully mapped.")
}

func (env *Env) getLogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	errLog, err := env.ErrorLog(i)
	if err != nil {
		log.Println(err)
		if errLog.ID == 0 {
			http.Error(w, err.Error(), 400)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}

	res := LogResponse{errLog}
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(res)
	fmt.Fprintf(w, string(j))
}
