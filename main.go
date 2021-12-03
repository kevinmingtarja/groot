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
	"strconv"
)

type Env struct {
	DB  *sql.DB
	Bot *tgbotapi.BotAPI
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	db, err := InitializeDB()
	if err != nil {
		log.Fatalln("Error establishing connection to database.")
	}
	bot, err := InitializeBot()
	if err != nil {
		log.Fatalln("Error initializing telegram bot.")
	}
	env := &Env{db, bot}

	r := mux.NewRouter()

	r.HandleFunc("/logs/url", env.getLogsByURLHandler).Methods("POST")
	r.HandleFunc("/logs", env.LogHandler).Methods("POST")
	r.HandleFunc("/logs/{id}", env.getLogHandler).Methods("GET")
	r.HandleFunc("/chat", env.setChatIDHandler).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)

		fmt.Fprintf(w, "Hello!")
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("Port is used.")
	}
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
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = env.SetChatID(&b)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
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
